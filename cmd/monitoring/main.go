package main

import (
	"net/http"

	"github.com/alexvelfr/Monitoring-go/internal/app/monitoring"
	"github.com/alexvelfr/Monitoring-go/internal/auth"
	"github.com/alexvelfr/Monitoring-go/internal/store"
	"github.com/gorilla/mux"
	"github.com/jasonlvhit/gocron"
)

func main() {
	defer store.DbStore.Close()

	router := mux.NewRouter().StrictSlash(true)
	core := router.PathPrefix("/core").Subrouter()
	core.Use(auth.BearerAuth)
	core.HandleFunc("/set-control-point", monitoring.IndexHandler).Methods(http.MethodPost)
	core.HandleFunc("/mailing", monitoring.MailingHandler).Methods(http.MethodPost)

	http.Handle("/", router)

	dispatcher := monitoring.CreateDispatcher()
	defer store.DbStore.Close()

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
