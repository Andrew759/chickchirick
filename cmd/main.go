package main

import (
	"chickChirick/cmd/config"
	"chickChirick/cmd/factory"
	"chickChirick/cmd/service"
	"chickChirick/internal/model/user"
)

// TODO: добавить логирование
func main() {
	factory.InitViper()

	appConfig := config.AppConfiguration{}.NewAppConfiguration()

	dbDecorator := service.InitORM(&appConfig.DatabaseConfig)
	defer dbDecorator.CloseDB()

	//TODO: не должно попасть в коммиты
	dbDecorator.GormInterface.AutoMigrate(&user.User{}, &user.Ban{}, &user.Meta{}, &user.Photo{}, &user.Property{})

	redisDecorator := service.InitRedis(appConfig.RedisConfig)
	defer redisDecorator.RedisClose()

	httpClient := factory.InitHttpClient()
	factory.BuildAndServe(dbDecorator, redisDecorator, httpClient)
}
