package monitoring

import (
	"fmt"
	"log"
	"time"

	"github.com/alexvelfr/Monitoring-go/internal/models"
	"github.com/alexvelfr/Monitoring-go/internal/store"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

//CheckBlocks - task which check cache
func CheckBlocks() {
	reglament := &models.Reglament{}
	for {
		var docs []models.Reglament = make([]models.Reglament, 10)
		err := store.DbStore.DB.Select(&docs, reglament.GetBlocksQuery())
		if err != nil {
			log.Println(err)
			continue
		}
		for _, doc := range docs {
			now := time.Now()
			downtime := int(now.Sub(doc.LastUpdated).Minutes())
			reglamentTime, reglamentName := getReglamentTime(doc)
			if downtime >= reglamentTime && doc.InReglament {
				doc.InReglament = false
				_, err := store.DbStore.DB.NamedExec(reglament.GetUpdateBlockQuery(), &doc)
				if err != nil {
					log.Println(err)
					continue
				}
				SendMassages(
					fmt.Sprintf("Блок %s вышел из регламента!\nВремя: %s\nРежим регламента: %s\nВремя регламента: %dмин.",
						doc.Block,
						time.Now().Format("01/02 15:04"),
						reglamentName,
						reglamentTime,
					))
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
