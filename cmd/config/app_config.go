package config

import (
	"chickChirick/cmd/config/dto"
	globalConfig "chickChirick/pkg/chirik_config"
	"github.com/spf13/viper"
)

type AppConfigurationInterface interface {
	NewAppConfiguration() AppConfiguration
}

type AppConfiguration struct {
	Environment string
	dto.DatabaseConfig
	RedisConfig
}

type RedisConfig struct {
	Host     string
	Port     int
	User     string
	Password string
}

func (c AppConfiguration) NewAppConfiguration() AppConfiguration {
	return AppConfiguration{
		Environment:    viper.GetString(globalConfig.Enviroment),
		DatabaseConfig: PrepareDatabaseConfig(),
		RedisConfig: RedisConfig{
			Host:     viper.GetString(globalConfig.RedisHost),
			Port:     viper.GetInt(globalConfig.RedisPort),
			User:     viper.GetString(globalConfig.RedisUser),
			Password: viper.GetString(globalConfig.RedisPassword),
		},
	}
}

func PrepareDatabaseConfig() dto.DatabaseConfig {
	dbc := dto.DatabaseConfig{}

	dbc.SetHost(viper.GetString(globalConfig.DbHost))
	dbc.SetPort(viper.GetInt(globalConfig.DbPort))
	dbc.SetName(viper.GetString(globalConfig.DbName))
	dbc.SetUser(viper.GetString(globalConfig.DbUser))
	dbc.SetPassword(viper.GetString(globalConfig.DbPass))
	dbc.SetTimezone(viper.GetString(globalConfig.DbTimezone))

	return dbc
}
