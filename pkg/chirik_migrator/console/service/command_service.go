package service

import (
	"chickChirick/pkg/chirik_migrator/console/config"
	"chickChirick/pkg/chirik_migrator/file"
	"fmt"
	"os"
)

type CommandService struct {
	MigratorService
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
		migratorEntities, err := file.ReadEntityDir(os.Args[2:])
		if err != nil {
			return err
		}

		err = cs.MigratorService.DoMigrate(migratorEntities)
		if err != nil {
			return err
		}

		return nil
	default:
		return fmt.Errorf("invalid command")
	}
}
