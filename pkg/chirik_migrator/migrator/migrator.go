package migrator

import (
	migratorDto "chickChirick/pkg/chirik_migrator/migrator/dto"
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

func run([][]Migrator) {

}

func (m Migrator) CreateTable(migratorEntities map[string][]migratorDto.MigratorInfo) error {
	for entityPath, migratorInfoList := range migratorEntities {
		for _, migratorInfo := range migratorInfoList {
			if !migratorInfo.MigratorEnabled {
				//TODO: возможно стоит предусмотреть тут выбрасывание ошибки
				continue
			}
			migratorInfo.
		}
	}
	//TODO: implement this
	return nil
}
