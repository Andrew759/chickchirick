package factory

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
)

func InitViper() {
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			panic(fmt.Errorf("config file not found: %w", err))
		}
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
}
