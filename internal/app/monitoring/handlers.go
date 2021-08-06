package monitoring

import (
	"log"
	"strings"

	"github.com/alexvelfr/Monitoring-go/internal/models"
	"github.com/alexvelfr/Monitoring-go/internal/store"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

//StartCommandHandler - process /start command
func StartCommandHandler(message *tgbotapi.Message) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Для продолжения регистрации, необходимо предоставить номер телефона!")
	msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButtonContact("Предоставить номер телефона!"),
	))

	var user models.User

	store.DbStore.DB.QueryRowx(user.SelectByTeledgramID(msg.ChatID)).StructScan(&user)

	if user.TelegramID == 0 {
		user.Name = message.Chat.FirstName + " " + message.Chat.LastName
		user.TelegramID = message.Chat.ID
		_, err := store.DbStore.DB.NamedExec(`INSERT INTO users (name, phone, telegram_id, send_messages, send_reports, is_admin) 
								VALUES (:name, :phone, :telegram_id, :send_messages, :send_reports, :is_admin)`, user)
		if err != nil {
			log.Println(err.Error())
		}
	}
	return msg
}

//ReportCommandHandler - process /repost command
func ReportCommandHandler(message *tgbotapi.Message) tgbotapi.MessageConfig {
	var user models.User
	var msg tgbotapi.MessageConfig
	store.DbStore.DB.QueryRowx(user.SelectByTeledgramID(message.Chat.ID)).StructScan(&user)
	if user.SendReports {
		msg = tgbotapi.NewMessage(message.Chat.ID, "Выберите вариант отчета")
		msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Статистика простоя за прошлый месяц", "#rep1")),
			tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Статистика простоя за текущий месяц", "#rep2")),
			tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Статистика простоя за вчера", "#rep3")),
			tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Статистика простоя за сегодня", "#rep4")),
		)
	} else {
		msg = tgbotapi.NewMessage(message.Chat.ID, "У вас нет права отчетов!")
	}

	return msg
}

// MessageAllHandler - process /mail command
func MessageAllHandler(message *tgbotapi.Message) tgbotapi.MessageConfig {
	var msg tgbotapi.MessageConfig
	var user models.User
	store.DbStore.DB.QueryRowx(user.SelectByTeledgramID(message.Chat.ID)).StructScan(&user)
	if user.IsAdmin {
		msg = tgbotapi.NewMessage(message.Chat.ID, "Рассылка отправлена!")
		text := strings.TrimSpace(strings.ReplaceAll(message.Text, "/mail", ""))
		if text != "" {
			go SendMassages(text)
		}
	} else {
		msg = tgbotapi.NewMessage(message.Chat.ID, "У вас нет права рассылок!")
	}
	return msg
}

//ContactHandler - process message with contact
func ContactHandler(message *tgbotapi.Message) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Вы успешно подписались на рассылку!")
	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)

	var user models.User

	store.DbStore.DB.QueryRowx(user.SelectByTeledgramID(msg.ChatID)).StructScan(&user)
	user.SendMessages = true
	user.Phone = message.Contact.PhoneNumber

	store.DbStore.DB.NamedExec(user.UpdateQuery(), user)
	return msg
}

//ReportCallbackHandler - process callback #rep1-4
func ReportCallbackHandler(callback *tgbotapi.CallbackQuery) tgbotapi.MessageConfig {
	var user models.User
	var msg tgbotapi.MessageConfig
	data := callback.Data

	store.DbStore.DB.QueryRowx(user.SelectByTeledgramID(callback.Message.Chat.ID)).StructScan(&user)

	if user.SendReports {
		reportText := getReportText(data)
		msg = tgbotapi.NewMessage(callback.Message.Chat.ID, reportText)
	} else {
		msg = tgbotapi.NewMessage(callback.Message.Chat.ID, "У вас нет права отчетов!")
	}
	msg.ParseMode = tgbotapi.ModeHTML
	return msg
}
