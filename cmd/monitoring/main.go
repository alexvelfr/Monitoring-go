package main

import (
	"net/http"
	"os"

	"github.com/alexvelfr/Monitoring-go/internal/app/monitoring"
	"github.com/alexvelfr/Monitoring-go/internal/auth"
	"github.com/alexvelfr/Monitoring-go/internal/store"
	"github.com/gorilla/mux"
	"github.com/jasonlvhit/gocron"
)

func createBaseConfig() {
	os.Mkdir("configs", os.ModePerm)
	file, _ := os.Create("configs/reglament.json")
	file.WriteString(`{
  "period": {
    "start": 8,
    "end": 22
  },
  "documents": [
    {
      "id": "Заявка на займ",
      "controlDay": 10,
      "controlNight": 30
    },
    {
      "id": "Пролонгация",
      "controlDay": 20,
      "controlNight": 60
    },
    {
      "id": "Выдача займа",
      "controlDay": 10,
      "controlNight": 60
    },
    {
      "id": "Погашения",
      "controlDay": 10,
      "controlNight": 60
    }
  ]
}
`)
	file.Close()
}

func main() {
	defer store.DbStore.Close()
	if _, err := os.Stat("configs"); os.IsNotExist(err) {
		createBaseConfig()
	}

	router := mux.NewRouter().StrictSlash(true)
	core := router.PathPrefix("/core").Subrouter()
	core.Use(auth.BearerAuth)
	core.HandleFunc("/set-control-point", monitoring.IndexHandler).Methods(http.MethodPost)
	core.HandleFunc("/mailing", monitoring.MailingHandler).Methods(http.MethodPost)

	http.Handle("/", router)

	dispatcher := monitoring.CreateDispatcher()

	go http.ListenAndServe("0.0.0.0:8444", nil)
	go monitoring.CheckBlocks()

	dispatcher.AddCommandHandler("start", monitoring.StartCommandHandler)

	dispatcher.AddCommandHandler("report", monitoring.ReportCommandHandler)

	dispatcher.AddCommandHandler("mail", monitoring.MessageAllHandler)

	dispatcher.AddCallbackHandler(`#rep\d`, monitoring.ReportCallbackHandler)

	dispatcher.RegisterContactHandler(monitoring.ContactHandler)

	gocron.Every(1).Day().At("09:00").DoSafely(monitoring.SendDailyReport)
	go func() {
		<-gocron.Start()
	}()

	for update := range monitoring.Bot.Updates {
		monitoring.Bot.Bot.Send(dispatcher.Dispatch(update))
	}
}
