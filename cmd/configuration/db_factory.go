package configuration

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"strings"
)

func InitDB(config DatabaseConfig) *gorm.DB {
	dsn := dns(config)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Errorf("db connect failed: %w", err))
	}
	return db
}

func dns(config DatabaseConfig) string {
	dsn := []string{
		"host=" + config.Host,
		"user=" + config.User,
		"password=" + config.Password,
		"dbname=" + config.Name,
		"port=" + string(rune(config.Port)),
		"TimeZone=" + config.Timezone,
	}
	return strings.Join(dsn, " ")
}
