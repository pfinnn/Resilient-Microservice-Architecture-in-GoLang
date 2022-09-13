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
	healthMysql "github.com/hellofresh/health-go/checks/mysql"
	"github.com/sirupsen/logrus"
)

type App struct {
	Router              *mux.Router
	Configuration       *configuration.Configuration
	IronPersistence    repositories.IronRepository
	SmelterService      service.SmelterService
}

func (app *App) Initialize() {
	var err error

	logrus.Debug("Initializing persistence(s).")
	app.IronPersistence, err = repositories.NewIronRepository(app.Configuration.MySQLConnection, app.Configuration.MySQLDatabase)
	errorcheck.CheckLogFatal(err)

	logrus.Debug("Setup healthchecks.")
	err = health.Register(health.Config{
		Name:      "mysql",
		Timeout:   time.Second * 3,
		SkipOnErr: false,
		Check: healthMysql.New(healthMysql.Config{
			DSN: app.Configuration.MySQLConnectionWithDatabase,
		}),
	})
	errorcheck.CheckLogFatal(err)

	logrus.Debug("Initializing HTTP-API.")
	app.SmelterService = service.NewSmelterService(
		app.IronPersistence,
		app.Configuration)
	handler := httpapi.NewHttpAPI(app.SmelterService)
	app.Router = mux.NewRouter()
	app.initializeHttpAPIRoutes(handler)
}

func (app *App) Cleanup() {
	app.IronPersistence.Cleanup()
}

func (app *App) initializeHttpAPIRoutes(handler *httpapi.HttpAPI) {
	eventsRouter := app.Router.PathPrefix("/iron").Subrouter()
	eventsRouter.Methods("GET").Path("").HandlerFunc(handler.GetIron)
	//eventsRouter.Methods("POST").Path("").HandlerFunc(handler.AddOrderHandler)

	app.Router.Handle("/status", health.Handler())
	app.Router.HandleFunc("/health", handler.HealthHandler).Methods("GET")
}

func (app *App) Serve(address string) {
	logrus.Debugf("Start serving HTTP-API on %s.", address)
	err := http.ListenAndServe(address, app.Router)
	errorcheck.CheckLogFatal(err)
}
