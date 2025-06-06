package main

import (
	"chickChirick/cmd/config"
	"chickChirick/cmd/factory"
	"chickChirick/cmd/service"
)

func main() {
	factory.InitViper()

	appConfig := config.AppConfiguration{}.NewAppConfiguration()

	dbDecorator := service.InitORM(&appConfig.DatabaseConfig)
	defer dbDecorator.CloseDB()

	redis := service.InitRedis(appConfig.RedisConfig)
	defer redis.RedisClose()

	//httpClient := factory.InitHttpClient()
	factory.InitServer()
}
