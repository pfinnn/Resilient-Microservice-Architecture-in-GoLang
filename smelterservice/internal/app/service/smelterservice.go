package service

import (
	"git.haw-hamburg.de/acm746/resilient-microservice/internal/app/configuration"
	"git.haw-hamburg.de/acm746/resilient-microservice/internal/app/domain/entities"
	"git.haw-hamburg.de/acm746/resilient-microservice/internal/app/domain/repositories"
	"github.com/sirupsen/logrus"
	"time"
)

type SmelterService interface {
	GetIron()(entities.Iron, error)
	SmeltOre()(entities.Iron, error)
}

type SmelterServiceImpl struct {
	ironRepository    repositories.IronRepository
	configuration      *configuration.Configuration
}

func NewSmelterService(
	ironRepository repositories.IronRepository,
	configuration *configuration.Configuration) SmelterService {

	return &SmelterServiceImpl{
		ironRepository:    ironRepository,
		configuration:      configuration,
	}
}

func (SmelterServiceImpl *SmelterServiceImpl) GetIron() (entities.Iron, error) {
	logrus.Debug("GetIron invoked")
	var err error
	iron, err := SmelterServiceImpl.SmeltOre()
	return iron, err
}

func (SmelterServiceImpl *SmelterServiceImpl) SmeltOre() (entities.Iron, error){
	logrus.Debug("SmeltOre invoked")
	var err error

	// Process Ore to Iron
	iron := entities.Iron{
		CreatedAt: time.Now(),
		Weight:    10,
		Quality:   10,
		Type:      "Iron",
	}
	return iron, err
}


