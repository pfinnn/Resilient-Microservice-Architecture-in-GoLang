package repositories

import (
	"git.haw-hamburg.de/acm746/resilient-microservice/internal/app/domain/entities"
)

type SwordRepository interface {
	GetSword() (entities.Sword, error)

	ExecuteCommand(interface{}) (interface{}, error)
	Cleanup()
}

func NewEventRepository(connection string) (SwordRepository, error) {
	return NewMongoDBLayer(connection)
}
