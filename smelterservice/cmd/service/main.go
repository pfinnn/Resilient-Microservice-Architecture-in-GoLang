package main

import (
	"git.haw-hamburg.de/acm746/resilient-microservice/internal/app"
	"git.haw-hamburg.de/acm746/resilient-microservice/internal/app/configuration"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Reading configuration
	application := app.App{}
	config := configuration.ReadConfiguration()
	application.Configuration = &config

	// Initialize logging
	logrus.SetOutput(os.Stdout)
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.DebugLevel)

	// Handle system signals for proper cleanup
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	application.Initialize()
	go application.Serve(application.Configuration.ServerAddress)

	// Wait until server signals quit
	select {
	case <-sigc:
		logrus.Debug("Cleaning up.")
		application.Cleanup()
		logrus.Debug("Exiting.")
	}
}
