package migrator

import (
	"chickChirick/cmd/service"
	migratorDto "chickChirick/pkg/chirik_migrator/migrator/provider"
	"fmt"
	"strconv"
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

	//TODO: переделать под подставляемые значения, а не включаемые на прямую в запрос
	resultSQL := "CREATE TABLE ? ("
	var sqlValues []string
	sqlValues = append(sqlValues, schema.Table)

	var fieldSQL string
	var fieldSQLList []string
	var fieldComment string
	var fieldCommentList []string

	fieldsCount := len(schema.Fields)
	for k, field := range schema.Fields {
		fieldSQL = field.Name + " " + field.DataType.String()
		if field.Size != 0 {
			fieldSQL += "(" + strconv.Itoa(field.Size) + ")"
		}
		fieldSQL += " "

		if field.PrimaryKey {
			fieldSQL += "PRIMARY KEY "
		}
		if field.AutoIncrement {
			fieldSQL += "AUTO_INCREMENT "
		}
		if field.HasDefaultValue {
			fieldSQL += "DEFAULT " + field.DefaultValue
		}
		if field.NotNull {
			fieldSQL += "NOT NULL "
		}
		if field.Unique {
			fieldSQL += "UNIQUE "
		}
		if field.Comment != "" {
			fieldComment = "comment on column " + schema.Name +
				"." + field.Name + " is '" + field.Comment + "';"

			fieldCommentList = append(fieldCommentList, fieldComment)
		}
		if k+1 < fieldsCount {
			fieldSQL += ","
		} else {
			fieldSQL += ";"
		}
		fieldSQLList = append(fieldSQLList, fieldSQL)
	}

	//TODO: удалить
	fmt.Println(resultSQL, fieldSQLList, fieldCommentList)

	return nil
}

//CREATE TABLE users (
//id INT AUTO_INCREMENT PRIMARY KEY,
//phone BIGINT,
//name VARCHAR(256),
//surname VARCHAR(256),
//password VARCHAR(1024)
//);
