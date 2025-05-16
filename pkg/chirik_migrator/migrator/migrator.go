package migrator

import (
	"chickChirick/cmd/service"
	migratorDto "chickChirick/pkg/chirik_migrator/migrator/provider"
)

//TODO: согласовать с интерфейсом

type Config struct {
	CreateIndexAfterCreateTable bool
	service.DBDecorator
}

type Migrator struct {
	Config
}

func (m Migrator) CreateTables(migratorEntities map[string][]migratorDto.MigratorInfo) error {
	//TODO: тут можно использовать entityPath вместо пустого вызова, стоит ли? Либо удалить
	for _, migratorInfoList := range migratorEntities {
		for _, migratorInfo := range migratorInfoList {
			return m.CreateTable(migratorInfo)
		}
	}
	return nil
}

func (m Migrator) CreateTable(migratorInfo migratorDto.MigratorInfo) error {
	if !migratorInfo.MigratorEnabled {
		//TODO: возможно стоит предусмотреть тут выбрасывание ошибки
		return nil
	}
	if migratorInfo.HasError() {
		//TODO доделать
	}

	return nil
}
