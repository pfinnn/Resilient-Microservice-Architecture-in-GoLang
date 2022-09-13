package httpapi

import (
	"encoding/json"
	"git.haw-hamburg.de/acm746/resilient-microservice/internal/app/domain/entities"
	"git.haw-hamburg.de/acm746/resilient-microservice/internal/app/service"
	"github.com/sirupsen/logrus"
	"net/http"
)

// API Setup

type HttpAPI struct {
	forgeService service.ForgeService
}

func NewHttpAPI(forgeService service.ForgeService) *HttpAPI {
	return &HttpAPI{
		forgeService: forgeService,
	}
}

func (eh *HttpAPI) HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (eh *HttpAPI) logRequest(r *http.Request) {
	logrus.WithFields(logrus.Fields{
		"addr":   r.RemoteAddr,
		"method": r.Method,
		"path":   r.URL.Path,
	}).Info("http request")
}

// API Methods

func (eh *HttpAPI) GetSwordHandler(w http.ResponseWriter, r *http.Request){
	eh.logRequest(r)

	var err error
	sword := entities.Sword{}
	sword, err = eh.forgeService.GetSword()
	if err != nil {
		logrus.Error("Error occured while trying to get a sword ")
		http.Error(w, `{"error": "Error occured while trying to get a sword"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json;charset=utf8")
	res, err := json.Marshal(&sword)
	if err == nil {
		_, err = w.Write(res)
	}
	if err != nil {
		logrus.Errorf("Error occured while trying to find all orders: %s", err.Error())
		http.Error(w, `{"error": "Error occured while trying to find all orders"}`, http.StatusInternalServerError)
		return
	}
}
