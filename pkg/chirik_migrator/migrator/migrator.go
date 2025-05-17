package migrator

import (
	"chickChirick/cmd/service"
	migratorDto "chickChirick/pkg/chirik_migrator/migrator/provider"
	"fmt"
)

//TODO: согласовать с интерфейсом

type Config struct {
	CreateIndexAfterCreateTable bool
	service.DBDecorator
}

type Migrator struct {
	Config
}

func (m Migrator) CreateTables(migratorEntities map[string][]migratorDto.MigratorInfo) error {
	for _, migratorInfoList := range migratorEntities {
		for _, migratorInfo := range migratorInfoList {
			return m.CreateTable(migratorInfo)
		}
	}
	return nil
}

func (m Migrator) CreateTable(migratorInfo migratorDto.MigratorInfo) error {
	if !migratorInfo.MigratorEnabled {
		return fmt.Errorf("can't process entity with disabled migrator: %s", migratorInfo.Schema.Name)
	}

	schema := &migratorInfo.Schema
	if migratorInfo.HasError() {
		return fmt.Errorf("can't process entity with errors at prepare stage: %s : %s",
			schema.Name,
			migratorInfo.ErrList,
		)
	}

	resultSQL := "CREATE TABLE ? ("

	var sqlValues []string

	sqlValues = append(sqlValues, schema.Table)

	var fieldSQL string
	var fieldSQLList []string

	fmt.Println(resultSQL, fieldSQLList)

	hasConfiguredPrimaryKey := schema.HasPrimaryKey()
	for key, field := range schema.Fields {
		fieldSQL += field.Name + " " + field.DataType.String() + " "
		//Первичный ключ устанавливается только один раз, если его по какой-то причине нет в конфиге
		if !hasConfiguredPrimaryKey && key == 0 {
			fieldSQL += "PRIMARY KEY "
		} else {
			if field.PrimaryKey {
				fieldSQL += "PRIMARY KEY "
			}
		}
	}

	return nil
}

//CREATE TABLE users (
//id INT AUTO_INCREMENT PRIMARY KEY,
//phone BIGINT,
//name VARCHAR(256),
//surname VARCHAR(256),
//password VARCHAR(1024)
//);
