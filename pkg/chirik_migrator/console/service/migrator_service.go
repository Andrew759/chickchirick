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
	errList := migratorService.CreateTables(migratorEntities)

	//TODO: временное решение на момент рефакторинга
	for _, err := range errList {
		return err
	}

	return nil
}
