package main_test

import (
	"git.haw-hamburg.de/acm746/resilient-microservice/internal/app/domain/entities"
	"github.com/magiconair/properties/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"

	"git.haw-hamburg.de/acm746/resilient-microservice/internal/app"
	"git.haw-hamburg.de/acm746/resilient-microservice/internal/app/configuration"
	"github.com/sirupsen/logrus"
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

func clearTable() {
	application.IronPersistence.Cleanup()
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

// Unit Tests

func TestGetIronSuccess(t *testing.T) {
	iron, err := application.SmelterService.GetIron()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, reflect.TypeOf(entities.Iron{}), reflect.TypeOf(iron))
}

func TestGetIronFailure(t *testing.T) {
	iron, err := application.SmelterService.GetIron()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, reflect.TypeOf(entities.Iron{}), reflect.TypeOf(iron))
}

func TestSmeltOre(t *testing.T) {
	iron, err := application.SmelterService.SmeltOre()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, reflect.TypeOf(entities.Iron{}), reflect.TypeOf(iron))
}