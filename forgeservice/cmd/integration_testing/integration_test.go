package integration_testing

import (
	"git.haw-hamburg.de/acm746/resilient-microservice/internal/app"
	"git.haw-hamburg.de/acm746/resilient-microservice/internal/app/configuration"
	"git.haw-hamburg.de/acm746/resilient-microservice/internal/app/domain/entities"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"os"
	"reflect"
	"testing"
)

var application app.App

func TestMain(m *testing.M) {
	// Reading configuration
	application = app.App{}
	config := configuration.ReadConfiguration(false)
	application.Configuration = &config

	// Initialize logging
	logrus.SetOutput(os.Stderr)
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.DebugLevel)

	application.Initialize()
	os.Exit(m.Run())
}

// Integration Tests

func Test_Integration_GetIronFromSmelterservice(t *testing.T){
	if testing.Short() {
		t.Skip("skipping Integration Test: TestGetIronFromSmelterservice")
	}
	iron, err := application.ForgeService.GetIronFromSmelterservice()
	if err != nil {
		t.Fail()
	}
	assert.Equal(t, reflect.TypeOf(entities.Iron{}), reflect.TypeOf(iron))
}

func Test_Integration_GetSword(t *testing.T){
	if testing.Short() {
		t.Skip("skipping Integration Test: TestGetSwords")
	}
	sword, err := application.ForgeService.GetSword()
	if err != nil {
		t.Fail()
	}
	assert.Equal(t, sword.Type, "Sword")
}

func Validate_Sword(sword entities.Sword) bool{
	return (sword.Type != "" && sword.Weight != 0)
}

func Test_Integration_GetSword_Success(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping Integration Test: TestGetSwords")
	}
	sword, err := application.ForgeService.GetSword()
	if err != nil {
		t.Fatal(err)
	}
	assert.True(t, Validate_Sword(sword))
}


