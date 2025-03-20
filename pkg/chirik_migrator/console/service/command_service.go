package service

import (
	"chickChirick/cmd/service"
	"chickChirick/pkg/chirik_migrator/console/config"
	"chickChirick/pkg/chirik_migrator/file"
	migratorDto "chickChirick/pkg/chirik_migrator/migrator/dto"
	"chickChirick/pkg/chirik_migrator/migrator/factory"
	"fmt"
	"os"
)

type CommandService struct {
	DBDecorator service.DBDecorator
}

func (cs CommandService) ParseInput() error {
	if len(os.Args) <= 1 {
		return fmt.Errorf("not enough input arguments")
	}
	command := os.Args[1]
	if len(command) <= 0 {
		return fmt.Errorf("empty command: %s", command)
	}

	switch command {
	case config.MigrateKey:
		//TODO: распараллелить?
		migratorEntities, err := file.ReadDir(os.Args[2:])
		if err != nil {
			return err
		}

		err = cs.doMigrate(migratorEntities)
		if err != nil {
			return err
		}

		return nil
	default:
		return fmt.Errorf("invalid command")
	}
}

// TODO: вынести отдельно как зависимость CommandService
func (cs CommandService) doMigrate(migratorEntities map[string][]migratorDto.MigratorInfo) error {
	migratorService := factory.InitMigrator(cs.DBDecorator, factory.WithCreateIndexAfterCreateTable())
	err := migratorService.CreateTables(migratorEntities)
	if err != nil {
		//TODO: реализовать
	}

	return nil
}
