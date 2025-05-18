package migrator

import (
	//TODO: тут при разнесении на микросервисы может быть проблема
	mainService "chickChirick/cmd/service"
	migratorDto "chickChirick/pkg/chirik_migrator/migrator/provider"
	"chickChirick/pkg/chirik_migrator/migrator/service"
	"fmt"
	"strconv"
)

//TODO: согласовать с интерфейсом

type Config struct {
	CreateIndexAfterCreateTable bool
	mainService.DBDecorator
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

	var sqlFieldList []string
	sqlFieldList = append(sqlFieldList, "CREATE TABLE %s (")

	var sqlValues []any
	sqlValues = append(sqlValues, schema.Table)

	var fieldCommentList []string
	var fieldCommentValues []string

	fieldsCount := len(schema.Fields)
	for k, field := range schema.Fields {
		sqlField := "%s %s"

		//TODO: При установке типа необходимость в дальнейшем вынесении отдельного провайдера для
		// postgres, т.к в разных БД реализация будет отличаться
		if field.AutoIncrement {
			sqlValues = append(sqlValues, field.Name, "BIGSERIAL")
		} else {
			sqlValues = append(sqlValues, field.Name, field.DataType.String())
		}

		if field.Size != 0 {
			sqlField += "(%s)"
			sqlValues = append(sqlValues, strconv.Itoa(field.Size))
		}

		if field.PrimaryKey {
			sqlField += " PRIMARY KEY"
		}
		if field.HasDefaultValue {
			sqlField += "DEFAULT %s"
			sqlValues = append(sqlValues, field.DefaultValue)
		}
		if field.NotNull {
			sqlField += " NOT NULL"
		}
		if field.Unique {
			sqlField += " UNIQUE"
		}
		if field.Comment != "" {
			fieldComment := "comment on column %s.%s is '%s';"
			fieldCommentValues = append(fieldCommentValues, schema.Name, field.Name, field.Comment)
			fieldCommentList = append(fieldCommentList, fieldComment)
		}
		if k+1 < fieldsCount {
			sqlField += ","
		} else {
			sqlField += ");"
		}
		sqlFieldList = append(sqlFieldList, sqlField)
	}

	//Предотвращение SQL инъекций по образу, как это делалось в PHP
	for k, v := range sqlValues {
		sqlValues[k] = service.Escape(v.(string))
	}

	//TODO: возможно имеет смысл сразу сетить всё в строку.
	var resultSQL string
	for _, sqlField := range sqlFieldList {
		resultSQL += sqlField
	}

	resultSQL = fmt.Sprintf(resultSQL, sqlValues...)

	result, err := m.NativeDB().Exec(resultSQL)

	//TODO: удалить. Можно вернуть результат и в отдельном сервисе записать в файл
	fmt.Println(result)

	return err
}

//CREATE TABLE users (
//id INT AUTO_INCREMENT PRIMARY KEY,
//phone BIGINT,
//name VARCHAR(256),
//surname VARCHAR(256),
//password VARCHAR(1024)
//);
