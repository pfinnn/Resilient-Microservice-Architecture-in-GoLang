package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
	"time"
)

func main() {
	logrus.SetOutput(os.Stderr)
	logrus.SetFormatter(&logrus.TextFormatter{})
	logrus.SetLevel(logrus.DebugLevel)

	if len(os.Args) != 3 {
		logrus.Fatalf("usage: %s <maxretrycount> <connection>", os.Args[0])
	}

	maxRetryCount, _ := strconv.Atoi(os.Args[1])
	var db *sql.DB
	var err error

	logrus.Debug("Waiting for MySQL to be ready.")
	for retryCount := 1; retryCount <= maxRetryCount; retryCount++ {
		db, err = sql.Open("mysql", os.Args[2])
		if err == nil {
			_, err := db.Exec("SELECT 1")
			if err == nil {
				logrus.Debug("MySQL is ready.")
				_ = db.Close()
				break
			}
		}

		logrus.Debugf("MySQL connection failed (retrying %d of %d every 3 seconds)", retryCount, maxRetryCount)
		time.Sleep(3 * time.Second)
	}
	if err != nil {
		logrus.Info("Could not establish connection to MySQL.")
		os.Exit(1);
	} else {
		os.Exit(0);
	}
}
