package monitoring

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/alexvelfr/Monitoring-go/internal/models"
	"github.com/alexvelfr/Monitoring-go/internal/store"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type block struct {
	ID           string `json:"id"`
	ControlDay   int    `json:"controlDay"`
	ControlNight int    `json:"controlNight"`
}
type config struct {
	Period struct {
		Start int `json:"start"`
		End   int `json:"end"`
	} `json:"period"`
	Documents []block `json:"documents"`
}

//CheckBlocks - task which check cache
func CheckBlocks() {
	configJSON, err := os.Open(filepath.Join("configs", "reglament.json"))

	if err != nil {
		log.Fatal(err)
	}

	var conf config
	decoder := json.NewDecoder(configJSON)
	decoder.Decode(&conf)
	configJSON.Close()

	for {
		for _, doc := range conf.Documents {
			reglament := &models.Reglament{}
			store.DbStore.DB.QueryRowx(reglament.GetBlockQuery(), doc.ID).StructScan(reglament)
			if reglament.ID == 0 {
				// Если блока не существует, создадим его
				reglament.Block = doc.ID
				reglament.InReglament = true
				reglament.LastUpdated = time.Now()
				store.DbStore.DB.NamedExec(reglament.GetCreateBlockQuery(), reglament)
				continue
			}
			now := time.Now()
			downtime := int(now.Sub(reglament.LastUpdated).Minutes())
			reglamentTime, reglamentName := getReglamentTime(&doc, &conf)
			if downtime >= reglamentTime && reglament.InReglament {
				reglament.InReglament = false
				store.DbStore.DB.NamedExec(reglament.GetUpdateBlockQuery(), reglament)
				SendMassages(fmt.Sprintf("Блок %s вышел из регламента!\nВремя: %s\nРежим регламента: %s\nВремя регламента: %dмин.", doc.ID, time.Now().Format("01/02 15:04"), reglamentName, reglamentTime))
			}
		}
		time.Sleep(time.Second * 15)
	}
}

// SendMassages - send messages for all users which has send_massages=true
func SendMassages(text string) {
	var recipients []models.User
	store.DbStore.DB.Select(&recipients, models.GetAllRecipients())
	for _, recipient := range recipients {
		msg := tgbotapi.NewMessage(recipient.TelegramID, text)
		Bot.Bot.Send(msg)
	}
}

//SendDailyReport - send daily report
func SendDailyReport() {
	report := getReportText("#rep3")
	var recipients []models.User
	store.DbStore.DB.Select(&recipients, models.GetAllReportRecipients())
	for _, recipient := range recipients {
		msg := tgbotapi.NewMessage(recipient.TelegramID, report)
		msg.ParseMode = tgbotapi.ModeHTML
		Bot.Bot.Send(msg)
	}
}
