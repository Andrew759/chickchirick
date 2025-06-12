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

	redisDecorator := service.InitRedis(appConfig.RedisConfig)
	defer redisDecorator.RedisClose()

	//httpClient := factory.InitHttpClient()
	factory.InitServer(dbDecorator, redisDecorator)
}
