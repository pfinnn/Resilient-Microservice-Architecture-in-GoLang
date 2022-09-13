package app

import (
	"net/http"
	"time"

	"git.haw-hamburg.de/acm746/resilient-microservice/internal/app/api/httpapi"
	"git.haw-hamburg.de/acm746/resilient-microservice/internal/app/configuration"
	"git.haw-hamburg.de/acm746/resilient-microservice/internal/app/domain/repositories"
	"git.haw-hamburg.de/acm746/resilient-microservice/internal/app/service"
	"git.haw-hamburg.de/acm746/resilient-microservice/internal/pkg/errorcheck"
	"github.com/gorilla/mux"
	"github.com/hellofresh/health-go"
	healthMongo "github.com/hellofresh/health-go/checks/mongo"
	"github.com/sirupsen/logrus"
)

type App struct {
	Router            	 	*mux.Router
	Configuration     		*configuration.Configuration
	SwordPersistence 		repositories.SwordRepository
	ForgeService    		service.ForgeService
}

func (app *App) Initialize() {
	var err error

	logrus.Debug("Initializing persistence(s).")
	app.SwordPersistence, err = repositories.NewEventRepository(app.Configuration.MongoConnection)
	errorcheck.CheckLogFatal(err)

	logrus.Debug("Setup healthchecks.")
	err = health.Register(health.Config{
		Name:      "mongo",
		Timeout:   time.Second * 3,
		SkipOnErr: false,
		Check: healthMongo.New(healthMongo.Config{
			DSN: app.Configuration.MongoConnection,
		}),
	})
	errorcheck.CheckLogFatal(err)
	logrus.Debug("Initializing HTTP-API.")
	app.ForgeService = service.NewForgeService(
		app.SwordPersistence,
		app.Configuration)
	handler := httpapi.NewHttpAPI(app.ForgeService)
	app.Router = mux.NewRouter()
	app.initializeHttpAPIRoutes(handler)
}

func (app *App) Cleanup() {
	app.SwordPersistence.Cleanup()
}

func (app *App) initializeHttpAPIRoutes(handler *httpapi.HttpAPI) {
	eventsRouter := app.Router.PathPrefix("/sword").Subrouter()
	eventsRouter.Methods("GET").Path("").HandlerFunc(handler.GetSwordHandler)

	app.Router.Handle("/status", health.Handler())
	app.Router.HandleFunc("/health", handler.HealthHandler).Methods("GET")
}

func (app *App) Serve(address string) {
	logrus.Debugf("Start serving HTTP-API on %s.", address)
	err := http.ListenAndServe(address, app.Router)
	errorcheck.CheckLogFatal(err)
}
