package httpapi

import (
	"encoding/json"
	"git.haw-hamburg.de/acm746/resilient-microservice/internal/app/service"
	"github.com/sirupsen/logrus"
	"net/http"
)

type HttpAPI struct {
	smelterService service.SmelterService
}

func NewHttpAPI(smelterService service.SmelterService) *HttpAPI {
	return &HttpAPI{
		smelterService: smelterService,
	}
}

func (eh *HttpAPI) HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (eh *HttpAPI) GetIron(w http.ResponseWriter, r *http.Request){
	eh.logRequest(r)

	iron, err := eh.smelterService.GetIron()

	if err != nil {
		logrus.Error("Error occured while trying to get iron")
		http.Error(w, `{"error": "Error occured while trying to get iron"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json;charset=utf8")
	res, err := json.Marshal(&iron)
	if err == nil {
		_, err = w.Write(res)
	}

}

func (eh *HttpAPI) logRequest(r *http.Request) {
	logrus.WithFields(logrus.Fields{
		"addr":   r.RemoteAddr,
		"method": r.Method,
		"path":   r.URL.Path,
	}).Info("http request")
}
