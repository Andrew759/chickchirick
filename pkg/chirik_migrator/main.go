package main

//TODO: Вынести в отдельный микросервис миграций?

import (
	"chickChirick/cmd/factory"
	appService "chickChirick/cmd/service"
	"chickChirick/internal/model/user"
	"chickChirick/pkg/chirik_migrator/config/dto"
	"chickChirick/pkg/chirik_migrator/console/service"
	"fmt"
)

func main() {
	loadConfig()
	service.ParseInput()

	dbConfig := dto.NewConfiguration()
	dbDecorator := appService.InitORM(&dbConfig)
	defer dbDecorator.CloseDB()

	//TODO: удалить мок
	err := dbDecorator.GormInterface.AutoMigrate(&user.User{})
	if err != nil {
		panic(fmt.Errorf("failed to migrate: %w", err))
	}
	fmt.Println(dbDecorator)
}

func loadConfig() {
	factory.InitViper()
	//TODO: поправить путь для продовой реализации
	factory.MergeConfigByFile("pkg/migration/migration.yaml")
}
