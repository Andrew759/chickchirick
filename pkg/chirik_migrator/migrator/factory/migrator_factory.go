package factory

import (
	"chickChirick/cmd/service"
	"chickChirick/pkg/chirik_migrator/migrator"
)

type migratorOptions struct { //Конфигурация структуры
	createIndexAfterCreateTable *bool
}

type MigratorOption func(options *migratorOptions)

func WithCreateIndexAfterCreateTable() MigratorOption { //Функция конфигурации,
	return func(mOptions *migratorOptions) {
		if mOptions.createIndexAfterCreateTable == nil {
			*mOptions.createIndexAfterCreateTable = true
		}
	}
}

func InitMigrator(dBDecorator service.DBDecorator, opts ...MigratorOption) migrator.Migrator {
	var mOptions migratorOptions
	for _, opt := range opts {
		opt(&mOptions)
	}

	return migrator.Migrator{
		migrator.Config{
			*mOptions.createIndexAfterCreateTable,
			dBDecorator,
		},
	}
}
