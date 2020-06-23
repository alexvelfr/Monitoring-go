package monitoring

import (
	"log"
	"os"
	s "strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

//Bot ...
var Bot *Monitoring = newMonitoring()

//Monitoring ...
type Monitoring struct {
	Bot     *tgbotapi.BotAPI
	Updates tgbotapi.UpdatesChannel
}

//New ...
func newMonitoring() *Monitoring {
	botToken, exist := os.LookupEnv("BOT_TOKEN")
	if !exist {
		log.Fatal("Token not found")
	}

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)
	url, exist := os.LookupEnv("url")
	if !s.HasSuffix(url, "/") {
		url = url + "/"
	}

	if exist {
		_, err = bot.SetWebhook(tgbotapi.NewWebhook(url + bot.Token))
		if err != nil {
			log.Fatal(err)
		}
	}

	info, err := bot.GetWebhookInfo()
	if err != nil {
		log.Fatal(err)
	}

	if info.LastErrorDate != 0 {
		log.Printf("Telegram callback failed: %s", info.LastErrorMessage)
	}

	return &Monitoring{
		Bot:     bot,
		Updates: bot.ListenForWebhook("/" + bot.Token),
	}
}
