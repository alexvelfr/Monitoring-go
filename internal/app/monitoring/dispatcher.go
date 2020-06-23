package monitoring

import (
	"log"
	"regexp"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type commandHanler struct {
	name    string
	handler func(message *tgbotapi.Message) tgbotapi.MessageConfig
}

type callbackHanler struct {
	data    *regexp.Regexp
	handler func(callback *tgbotapi.CallbackQuery) tgbotapi.MessageConfig
}

//Dispatcher ...
type Dispatcher struct {
	commandHanlers  []commandHanler
	callbackHanlers []callbackHanler
	contactHandler  func(message *tgbotapi.Message) tgbotapi.MessageConfig
}

//CreateDispatcher ...
func CreateDispatcher() *Dispatcher {
	return &Dispatcher{
		contactHandler: func(message *tgbotapi.Message) tgbotapi.MessageConfig {
			return tgbotapi.MessageConfig{}
		},
	}
}

//Dispatch incoming update
func (d *Dispatcher) Dispatch(update tgbotapi.Update) tgbotapi.MessageConfig {
	if update.Message != nil && update.Message.IsCommand() {
		return d.processCommand(update.Message)
	}
	if update.CallbackQuery != nil {
		return d.processCallback(update.CallbackQuery)
	}
	if update.Message.Contact != nil {
		return d.contactHandler(update.Message)
	}
	return tgbotapi.MessageConfig{}
}

//AddCallbackHandler - add new func handler for callback data
func (d *Dispatcher) AddCallbackHandler(data string, handl func(callback *tgbotapi.CallbackQuery) tgbotapi.MessageConfig) {
	re, err := regexp.Compile(data)
	if err != nil {
		log.Fatal("Cannot compiler regexp " + data)
	}
	ch := callbackHanler{
		data:    re,
		handler: handl,
	}
	callbackHanlers := append(d.callbackHanlers, ch)
	d.callbackHanlers = callbackHanlers
}

//AddCommandHandler - add new func handler for command
func (d *Dispatcher) AddCommandHandler(name string, handl func(message *tgbotapi.Message) tgbotapi.MessageConfig) {
	ch := commandHanler{
		name:    name,
		handler: handl,
	}
	commandHanlers := append(d.commandHanlers, ch)
	d.commandHanlers = commandHanlers
}

//RegisterContactHandler - register func handler for contact recive
func (d *Dispatcher) RegisterContactHandler(handl func(message *tgbotapi.Message) tgbotapi.MessageConfig) {
	d.contactHandler = handl
}

func (d *Dispatcher) processCommand(message *tgbotapi.Message) tgbotapi.MessageConfig {
	for _, command := range d.commandHanlers {
		if command.name == message.Command() {
			return command.handler(message)
		}
	}
	return tgbotapi.MessageConfig{}
}

func (d *Dispatcher) processCallback(callback *tgbotapi.CallbackQuery) tgbotapi.MessageConfig {
	for _, command := range d.callbackHanlers {
		if command.data.MatchString(callback.Data) {
			return command.handler(callback)
		}
	}
	return tgbotapi.MessageConfig{}
}
