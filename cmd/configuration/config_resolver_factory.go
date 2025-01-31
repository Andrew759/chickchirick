package configuration

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
)

// InitViper TODO: рассмотреть инициализацию
func InitViper() {
	viper.SetConfigName("config")
	viper.SetConfigType("env")        // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("../../.env") // path to look for the config file in
	viper.AddConfigPath("../.env")    // optionally look for config in the working directory
	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			panic(fmt.Errorf("config file not found: %w", err))
		}
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
}
