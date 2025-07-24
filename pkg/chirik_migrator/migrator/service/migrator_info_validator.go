package service

import (
	"chickChirick/pkg/chirik_migrator/db_schema"
	migratorDto "chickChirick/pkg/chirik_migrator/migrator/provider"
	"errors"
	"fmt"
)

func ValidateMInfo(migratorInfo migratorDto.MigratorInfo) error {
	if !migratorInfo.MigratorEnabled {
		return fmt.Errorf("can't process entity with disabled migrator: %s", migratorInfo.Schema.Name)
	}

	if migratorInfo.HasCriticalError() {
		return fmt.Errorf("can't process entity with errors at prepare stage: %s : %s",
			migratorInfo.Schema.Name,
			migratorInfo.ErrList,
		)
	}

	if migratorInfo.HasInfoError() {
		var fieldTypeError *db_schema.FieldTypeError
		for _, err := range migratorInfo.ErrList {
			if errors.As(err, &fieldTypeError) {
				//TODO: Implement this
			}
		}
	}

	return nil
}
