package dto

import mainService "chickChirick/cmd/service"

type MigratorDiContainer struct {
	MConfig
	mainService.DBDecorator
}

type MConfig struct {
	//TOOD: CreateIndexAfterCreateTable сейчас не используется. Проверить необходимость
	CreateIndexAfterCreateTable bool
	MigrationFilesPath          string
	EnableTableNamespace        bool
	FixtureCount                int
	FixturePrefix               string
	FixtureNilColumns           map[string]struct{}
}
