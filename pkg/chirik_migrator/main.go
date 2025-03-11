package main

//TODO: Вынести в отдельный микросервис миграций?

import (
	"chickChirick/cmd/factory"
	appService "chickChirick/cmd/service"
	"chickChirick/pkg/chirik_migrator/config/dto"
	"chickChirick/pkg/chirik_migrator/console/service"
	"fmt"
)

func main() {
	loadConfig()
	err := service.ParseInput()
	if err != nil {
		panic(fmt.Errorf("failed to parse input: %w", err))
	}

	dbConfig := dto.NewConfiguration()
	dbDecorator := appService.InitORM(&dbConfig)
	defer dbDecorator.CloseDB()
}

func loadConfig() {
	factory.InitViper()
	//TODO: поправить путь для продовой реализации
	factory.MergeConfigByFile("pkg/chirik_migrator/migration.yaml")
}
