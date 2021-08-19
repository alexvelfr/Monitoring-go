package monitoring

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type document struct {
	Name string `json:"document"`
}

type requestMailing struct {
	Params struct {
		Message string
		Service struct{ Status string }
	} `json:"params"`
}

//IndexHandler - index handler
func IndexHandler(w http.ResponseWriter, r *http.Request) {

	bytes, _ := ioutil.ReadAll(r.Body)
	res := make(map[string]string)
	data := document{}
	res["succes"] = "ok"

	if err := json.Unmarshal(bytes, &data); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error format"))
		return
	}
	processDocument(data)

	resBt, err := json.Marshal(res)
	if err != nil {
		log.Print(err.Error())
	}
	w.Write(resBt)
}

//MailingHandler - mailing handler
func MailingHandler(w http.ResponseWriter, r *http.Request) {
	bytes, _ := ioutil.ReadAll(r.Body)
	res := make(map[string]string)
	data := requestMailing{}
	res["succes"] = "ok"

	if err := json.Unmarshal(bytes, &data); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error format"))
		return
	}

	if data.Params.Service.Status != "" {
		processServiceMessage(&data)
	}
	SendMassages(data.Params.Message)

	resBt, err := json.Marshal(res)
	if err != nil {
		log.Print(err.Error())
	}
	w.Write(resBt)
}
