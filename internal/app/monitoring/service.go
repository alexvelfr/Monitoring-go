package monitoring

import (
	"fmt"
	"log"
	"time"

	"github.com/alexvelfr/Monitoring-go/internal/models"
	"github.com/alexvelfr/Monitoring-go/internal/store"
)

func processDocument(doc document) {
	var reglament models.Reglament

	store.DbStore.DB.QueryRowx(reglament.GetBlockByCodeQuery(), doc.Name).StructScan(&reglament)

	if reglament.ID == 0 {
		reglament.Block = doc.Name
		reglament.InReglament = true
		reglament.Active = true
		reglament.LastUpdated = time.Now()
		reglament.Code = doc.Name
		reglament.DayHour = 8
		reglament.NightHour = 22
		reglament.ReglamentDayTime = 10
		reglament.ReglamentNightTime = 60
		_, err := store.DbStore.DB.NamedExec(reglament.GetCreateBlockQuery(), &reglament)
		if err != nil {
			log.Println(err)
		}
		return
	}
	if !reglament.InReglament {
		go processReturnInReglament(reglament)
	}
	reglament.LastUpdated = time.Now()
	reglament.InReglament = true
	store.DbStore.DB.NamedExec(reglament.GetUpdateBlockQuery(), reglament)
}

func processReturnInReglament(reglament models.Reglament) {
	loc, _ := time.LoadLocation("Europe/Kiev")
	now := time.Now().In(loc)
	downtime := int(now.Sub(reglament.LastUpdated.In(loc)).Seconds())
	statistic := &models.Statistic{
		Block:             reglament.Block,
		OutOfReglament:    reglament.LastUpdated.In(loc),
		ReturnInReglament: now,
		Downtime:          downtime,
	}
	store.DbStore.DB.NamedExec(statistic.CreateNewStatiscitQuery(), statistic)
	SendMassages(fmt.Sprintf("Вернулся к работе блок %s!\nВремя: %s\nВремя простоя: %dмин.", reglament.Block, time.Now().Format("01/02 15:04"), int(downtime/60)))
}

func getReglamentTime(bl models.Reglament) (int, string) {
	now := time.Now().Hour()
	if now >= bl.DayHour && now < bl.NightHour {
		return bl.ReglamentDayTime, "День"
	}
	return bl.ReglamentNightTime, "Ночь"
}

func processServiceMessage(data *requestMailing) {
	const BLOCK1C string = "1C"
	if data.Params.Service.Status == "down" {
		var reglament models.Reglament
		store.DbStore.DB.QueryRowx(reglament.GetBlockQuery(), BLOCK1C).StructScan(&reglament)
		reglament.Block = BLOCK1C
		reglament.InReglament = false
		reglament.LastUpdated = time.Now()
		store.DbStore.DB.NamedExec(reglament.GetUpdateBlockQuery(), &reglament)
	}
	if data.Params.Service.Status == "up" {
		var reglament models.Reglament
		store.DbStore.DB.QueryRowx(reglament.GetBlockQuery(), BLOCK1C).StructScan(&reglament)
		if reglament.ID == 0 || reglament.InReglament {
			return
		}

		downtime := int(time.Since(reglament.LastUpdated).Seconds())
		//Добавим сообщение о времени простоя сервиса
		data.Params.Message = data.Params.Message + fmt.Sprintf("\nВремя простоя: %d мин.", int(downtime/60))
		statistic := &models.Statistic{
			Block:             BLOCK1C,
			OutOfReglament:    reglament.LastUpdated,
			ReturnInReglament: time.Now(),
			Downtime:          downtime,
		}
		store.DbStore.DB.NamedExec(statistic.CreateNewStatiscitQuery(), statistic)
		reglament.InReglament = true
		reglament.LastUpdated = time.Now()
		store.DbStore.DB.NamedExec(reglament.GetUpdateBlockQuery(), &reglament)
	}
}

func getReportText(repNumber string) string {
	statistic := &models.Statistic{}
	query, param1, param2 := statistic.GetReportQuery(repNumber)
	title := statistic.GetReportTitle(repNumber)
	reportBody := ""
	var report []struct {
		Block    string `db:"block"`
		Downtime int    `db:"downtime"`
		Counts   int    `db:"counts"`
	}
	store.DbStore.DB.Select(&report, query, param1, param2)
	for _, rep := range report {
		template := `
		<b>Блок: %s</b>
		Выходов из регламента: %d
		Время простоя в часах:  %d
		Время простоя в минутах: %d
		=============================
		`
		reportBody += fmt.Sprintf(template, rep.Block, rep.Counts, int(rep.Downtime/60), rep.Downtime)
	}
	return title + reportBody
}
