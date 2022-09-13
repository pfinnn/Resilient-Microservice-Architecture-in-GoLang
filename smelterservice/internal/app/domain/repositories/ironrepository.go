package repositories

import (
	"github.com/jinzhu/gorm"
)

type IronRepository interface {

	GetSQLConnection() *gorm.DB
	Cleanup()
}

func NewIronRepository(connection string, databaseName string) (IronRepository, error) {
	return MySQLBLayer(connection, databaseName)
}
