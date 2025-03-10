package main

//TODO: Вынести в отдельный микросервис миграций?

import (
	"chickChirick/cmd/factory"
	appService "chickChirick/cmd/service"
	"chickChirick/internal/model/user"
	"chickChirick/pkg/migration/console/config"
	"chickChirick/pkg/migration/console/service"
	"fmt"
)

func main() {
	loadConfig()
	service.ParseInput()

	dbConfig := config.NewConfiguration()
	dbDecorator := appService.InitORM(&dbConfig)
	defer dbDecorator.CloseDB()

	err := dbDecorator.ORMInterface.AutoMigrate(&user.User{})
	if err != nil {
		panic(fmt.Errorf("failed to migrate: %w", err))
	}

	fmt.Println(dbDecorator)
}

func loadConfig() {
	factory.InitViper()
	factory.MergeConfigByFile("migration/console/migration.yaml")
}
