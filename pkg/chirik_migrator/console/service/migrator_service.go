package service

import (
	"chickChirick/cmd/service"
	"chickChirick/pkg/chirik_migrator/migrator/factory"
	migratorDto "chickChirick/pkg/chirik_migrator/migrator/provider"
)

type MigratorService struct {
	DBDecorator service.DBDecorator
}

func (ms MigratorService) DoMigrate(migratorEntities map[string][]migratorDto.MigratorInfo) error {
	migratorService := factory.InitMigrator(ms.DBDecorator, factory.WithCreateIndexAfterCreateTable())
	return migratorService.CreateTables(migratorEntities)
}
