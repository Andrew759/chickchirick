package migrator

import (
	"database/sql"
	"gorm.io/gorm"
)

type DBDecorator struct {
	ORMInterface    *gorm.DB
	NativeInterface *sql.DB
}
type Config struct {
	CreateIndexAfterCreateTable bool
	DBDecorator
}

type Migrator struct {
	Config
}

func (m Migrator) CreateTable() error {
	//TODO: implement this
	return nil
}
