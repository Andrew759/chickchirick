package service

import (
	"chickChirick/cmd/configuration"
	"database/sql"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"strconv"
	"strings"
)

type DBDecorator struct {
	ORMInterface    *gorm.DB
	NativeInterface *sql.DB
}

func InitAndPrepareORM(config configuration.DatabaseConfig) DBDecorator {
	dsn := dsn(config)

	ORM, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Errorf("db connect failed: %w", err))
	}

	//TODO: предусмотреть обработку ошибок?
	nativeDB, _ := ORM.DB()
	dbd := DBDecorator{
		ORMInterface:    ORM,
		NativeInterface: nativeDB,
	}

	dbd.DeferDBClose()

	return dbd
}

func dsn(config configuration.DatabaseConfig) string {
	dsn := []string{
		"host=" + config.Host,
		"user=" + config.User,
		"password=" + config.Password,
		"dbname=" + config.Name,
		"port=" + strconv.Itoa(config.Port),
	}
	if config.Timezone != "" {
		dsn = append(dsn, "TimeZone="+config.Timezone)
	}

	return strings.Join(dsn, " ")
}

func (dbd DBDecorator) DeferDBClose() {
	defer func(NativeInterface *sql.DB) {
		err := NativeInterface.Close()
		if err != nil {
			panic(fmt.Errorf("db close error: %w", err))
		}
	}(dbd.NativeInterface)
}
