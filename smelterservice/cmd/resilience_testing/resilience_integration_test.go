package resilience_testing

import (
	"git.haw-hamburg.de/acm746/resilient-microservice/internal/app"
	"git.haw-hamburg.de/acm746/resilient-microservice/internal/app/configuration"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var application app.App



func TestMain(m *testing.M) {
	// Reading configuration
	application = app.App{}
	config := configuration.ReadConfiguration()
	application.Configuration = &config

	// Initialize logging
	logrus.SetOutput(os.Stderr)
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.DebugLevel)

	application.Initialize()
	os.Exit(m.Run())
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	application.Router.ServeHTTP(rr, req)

	return rr
}
func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}
