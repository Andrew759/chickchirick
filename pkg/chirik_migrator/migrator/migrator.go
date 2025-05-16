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
	//TODO: тут можно использовать entityPath вместо пустого вызова, стоит ли? Либо удалить
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

	fmt.Sprintf(resultSQL, fieldSQLList)

	hasConfiguredPrimaryKey := schema.HasPrimaryKey()
	for _, field := range schema.Fields {
		fieldSQL += field.Name + " " + field.DataType.String()
		if !hasConfiguredPrimaryKey {

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
