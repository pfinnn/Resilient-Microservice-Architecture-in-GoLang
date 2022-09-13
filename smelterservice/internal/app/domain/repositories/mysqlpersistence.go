package repositories

import (
	"errors"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/sirupsen/logrus"
	"git.haw-hamburg.de/acm746/resilient-microservice/internal/app/domain/entities"
	"time"
)

type MySQLDBLayer struct {
	db *gorm.DB
}

func (mySQLDBLayer *MySQLDBLayer) GetIron() (entities.Iron, error) {
	panic("GetIron is not supposed to access DB, dont implement me and dont call me")
}

var (
	ErrIdMustNotBeSet   = errors.New("Id must not be set")
)

func MySQLBLayer(connection string, databaseName string) (IronRepository, error) {
	const maxRetryCount = 15
	var db *gorm.DB
	var err error

	logrus.Debug("Connecting to MySQL.")
	for retryCount := 1; retryCount <= maxRetryCount; retryCount++ {
		db, err = gorm.Open("mysql", connection + "?parseTime=true")
		if err == nil {
			err = db.Exec("SELECT 1").Error
			if err == nil {
				logrus.Debug("MySQL connection successfully established.")
				break
			}
		}

		logrus.Debugf("MySQL connection failed with error (retrying %d of %d): %s", retryCount, maxRetryCount, err)
		time.Sleep(3 * time.Second)
	}
	if err != nil {
		logrus.Fatal("Could not establish connection to MySQL.")
	}

	err = db.Exec("CREATE DATABASE IF NOT EXISTS " + databaseName + ";").Error
	if err == nil {
		logrus.Debugf("Database '%s' created.", databaseName)
	} else {
		logrus.Debugf("Re-using existing database '%s'.", databaseName)
	}

	err = db.Close()
	if err != nil {
		logrus.Fatal("Could not close MySQL-connection.")
	}
	db, err = gorm.Open("mysql", connection + databaseName + "?parseTime=true")
	if err != nil {
		logrus.Fatalf("Could not re-establish connection to MySQL now using database %s.", databaseName)
	}
	if db != nil {
		if db.Error != nil {
			logrus.Fatalf("Could not re-establish connection to MySQL now using database %s.", databaseName)
		}
	}

	db = db.Set("gorm:auto_preload", true)

	err = db.AutoMigrate(&entities.Iron{}).Error
	if err != nil {
		logrus.Fatalf("Unable to migrate database %s.", databaseName)
	}
	logrus.Debugf("Database '%s' migrated.", databaseName)

	db.LogMode(false)

	return &MySQLDBLayer{
		db: db,
	}, err
}


func (mySQLDBLayer *MySQLDBLayer) DeleteAllData() {
	mySQLDBLayer.db.Unscoped().Where("").Delete(entities.Iron{})
}

func (mySQLDBLayer *MySQLDBLayer) GetSQLConnection() *gorm.DB {
	return mySQLDBLayer.db
}

func (mySQLDBLayer *MySQLDBLayer) Cleanup() {
	logrus.Debug("Closing MySQL connection.")
	_ = mySQLDBLayer.db.Close()
}