package repositories

import (
	"git.haw-hamburg.de/acm746/resilient-microservice/internal/app/domain/entities"
	"git.haw-hamburg.de/acm746/resilient-microservice/internal/pkg/errorcheck"
	"github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2"
	"time"
)

type MongoDBLayer struct {
	session *mgo.Session
}

func NewMongoDBLayer(connection string) (SwordRepository, error) {
	logrus.Debugf("Connecting to mongodb on %s.", connection)
	session, err := mgo.DialWithTimeout(connection, 30 * time.Second)
	errorcheck.CheckLogFatal(err)
	logrus.Debug("Connection to MongoDB successfully established.")

	return &MongoDBLayer{
		session: session,
	}, err
}

const (
	MONGODB_SWORDSCOLLECTION = "swords"
)

func (mgoLayer *MongoDBLayer) GetSword() (entities.Sword, error) {
	session := mgoLayer.getFreshSession()
	defer session.Close()

	sword := entities.Sword{}
	err := session.DB("").C(MONGODB_SWORDSCOLLECTION).Find(nil).One(&sword)
	return sword, err
}

func (mgoLayer *MongoDBLayer) getFreshSession() *mgo.Session {
	return mgoLayer.session.Copy()
}

func (mgoLayer *MongoDBLayer) ExecuteCommand(cmd interface{}) (interface{}, error) {
	session := mgoLayer.getFreshSession()
	defer session.Close()

	var result interface{}
	err := session.DB("").Run(cmd, result)
	return result, err
}

func (mgoLayer *MongoDBLayer) Cleanup() {
	logrus.Debug("Closing Mongo connection.")
	_ = mgoLayer.session.DB("").DropDatabase()
}