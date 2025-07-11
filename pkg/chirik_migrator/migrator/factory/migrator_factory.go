package factory

import (
	"chickChirick/cmd/service"
	"chickChirick/pkg/chirik_migrator/config"
	"chickChirick/pkg/chirik_migrator/migrator"
	"github.com/spf13/viper"
)

type migratorOptions struct { //Конфигурация структуры
	createIndexAfterCreateTable bool
}

type MigratorOption func(options *migratorOptions)

func WithCreateIndexAfterCreateTable() MigratorOption { //Функция конфигурации,
	return func(mOptions *migratorOptions) {
		mOptions.createIndexAfterCreateTable = true
	}
}

func InitMigrator(dBDecorator service.DBDecorator, opts ...MigratorOption) migrator.Migrator {
	var mOptions migratorOptions
	for _, opt := range opts {
		opt(&mOptions)
	}

	return migrator.Migrator{
		Config: migrator.Config{
			CreateIndexAfterCreateTable: mOptions.createIndexAfterCreateTable,
			DBDecorator:                 dBDecorator,
			MigrationFilesPath:          viper.GetString(config.MigrationPath),
			EnableTableNamespace:        viper.GetBool(config.EnableTableNamespace),
			FixtureCount:                viper.GetInt(config.FixtureCount),
			FixturePrefix:               viper.GetString(config.FixturePrefix),
		},
	}
}
