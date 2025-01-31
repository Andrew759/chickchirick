package configuration

import (
	globalConfig "chickChirick/pkg/configuration"
	"github.com/spf13/viper"
)

// AppConfiguration TODO: возможно это надо вынести в модули
type AppConfiguration struct {
	Environment string
	DatabaseConfig
	RedisConfig
}

type DatabaseConfig struct {
	Host          string
	Port          int
	Name          string
	User          string
	Password      string
	Timezone      string
	MigrationPath string
}

type RedisConfig struct {
	Host string
	Port int
	Name string
}

func NewConfiguration() AppConfiguration {
	return AppConfiguration{
		viper.GetString(globalConfig.Enviroment),
		DatabaseConfig{
			Host:          viper.GetString(globalConfig.DbHost),
			Port:          viper.GetInt(globalConfig.DbPort),
			Name:          viper.GetString(globalConfig.DbName),
			User:          viper.GetString(globalConfig.DbUser),
			Password:      viper.GetString(globalConfig.DbPass),
			Timezone:      viper.GetString(globalConfig.DbTimezone),
			MigrationPath: viper.GetString(globalConfig.DbMigrationPatch),
		},
		RedisConfig{
			Host: viper.GetString(globalConfig.RedisHost),
			Port: viper.GetInt(globalConfig.RedisPort),
			Name: viper.GetString(globalConfig.RedisDb),
		},
	}
}
