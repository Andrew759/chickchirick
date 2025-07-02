package migrator

import (
	mainService "chickChirick/cmd/service"
	"chickChirick/pkg/chirik_migrator/db_schema"
	"chickChirick/pkg/chirik_migrator/file"
	"chickChirick/pkg/chirik_migrator/migrator/dto"
	migratorDto "chickChirick/pkg/chirik_migrator/migrator/provider"
	"chickChirick/pkg/chirik_migrator/migrator/service"
	"errors"
	"fmt"
	"strconv"
)

//TODO: согласовать с интерфейсом

type Config struct {
	//TOOD: CreateIndexAfterCreateTable сейчас не используется. Проверить необходимость
	CreateIndexAfterCreateTable bool
	mainService.DBDecorator
	MigrationFilesPath    string
	EnableTableNamespace  bool
	EnableCreatedAtColumn bool
	EnableUpdatedAtColumn bool
	EnableDeleteAtColumn  bool
	EnableFixtures        bool
}

type Migrator struct {
	Config
}

func (m Migrator) CreateTables(migratorEntities map[string][]migratorDto.MigratorInfo) error {
	for _, migratorInfoList := range migratorEntities {
		for _, migratorInfo := range migratorInfoList {
			err := m.CreateTable(migratorInfo)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (m Migrator) CreateTable(migratorInfo migratorDto.MigratorInfo) error {
	err := m.validateMInfo(migratorInfo)
	if err != nil {
		return err
	}

	sqlMeta := m.processSchemaFields(migratorInfo)
	m.addMigratorFields(&sqlMeta)

	//Предотвращение SQL инъекций по образу, как это делалось в PHP
	for k, v := range sqlMeta.SqlValues {
		sqlMeta.SqlValues[k] = service.Escape(v.(string))
	}

	var resultSQL string
	for _, sqlField := range sqlMeta.SqlFieldList {
		resultSQL += sqlField
	}

	resultSQL = fmt.Sprintf(resultSQL, sqlMeta.SqlValues...)

	_, err = m.NativeDB().Exec(resultSQL)
	if err != nil {
		return err
	}

	err = file.WriteSQLToFile(resultSQL, sqlMeta.TableName, m.MigrationFilesPath)
	if err != nil {
		return err
	}

	if m.EnableFixtures {
		m.insertFixtures(sqlMeta)
	}

	return err
}

func (m Migrator) validateMInfo(migratorInfo migratorDto.MigratorInfo) error {
	if !migratorInfo.MigratorEnabled {
		return fmt.Errorf("can't process entity with disabled migrator: %s", migratorInfo.Schema.Name)
	}

	if migratorInfo.HasCriticalError() {
		return fmt.Errorf("can't process entity with errors at prepare stage: %s : %s",
			migratorInfo.Schema.Name,
			migratorInfo.ErrList,
		)
	}

	if migratorInfo.HasInfoError() {
		var fieldTypeError *db_schema.FieldTypeError
		for _, err := range migratorInfo.ErrList {
			if errors.As(err, &fieldTypeError) {
				//TODO: Implement this
			}
		}
	}

	return nil
}

func (m Migrator) processSchemaFields(migratorInfo migratorDto.MigratorInfo) dto.Meta {
	var sqlFieldList []string
	sqlFieldList = append(sqlFieldList, "CREATE TABLE IF NOT EXISTS %s \n(")

	schema := &migratorInfo.Schema
	fullTableName := schema.Table
	if m.Config.EnableTableNamespace {
		fullTableName = migratorInfo.EntityNamespace + "_" + schema.Table
	}

	var sqlValues []any
	sqlValues = append(sqlValues, fullTableName)

	var fieldCommentList []string
	var fieldCommentValues []string
	var fieldMetas []dto.FieldMeta
	fieldCount := len(schema.Fields)
	processedCount := 0

	for _, field := range schema.Fields {
		fieldType := field.DataType.String()
		//Пропуск полей без типа
		if fieldType == "" {
			continue
		}

		sqlField := "\n %s %s"

		//TODO: При установке типа необходимость в дальнейшем вынесении отдельного провайдера для
		// postgres, т.к в разных БД реализация будет отличаться
		if field.AutoIncrement {
			fieldType = string(db_schema.BigSerial)
		}
		sqlValues = append(sqlValues, field.Name, fieldType)

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

		processedCount++
		if processedCount < fieldCount {
			sqlField += ","
		}

		sqlFieldList = append(sqlFieldList, sqlField)

		fieldMetas = append(fieldMetas, dto.FieldMeta{
			Name: field.Name,
			Type: fieldType,
		})
	}

	return dto.Meta{
		TableName:          fullTableName,
		SqlFieldList:       sqlFieldList,
		SqlValues:          sqlValues,
		FieldCommentList:   fieldCommentList,
		FieldCommentValues: fieldCommentValues,
		FieldMetas:         fieldMetas,
		FieldCount:         fieldCount,
	}
}

func (m Migrator) addMigratorFields(sqlMeta *dto.Meta) {
	if m.EnableCreatedAtColumn {
		sqlField := ",\n created_at " +
			db_schema.TimestampWithTimezone.String() + " " +
			db_schema.Null.String()
		sqlMeta.SqlFieldList = append(sqlMeta.SqlFieldList, sqlField)
	}
	if m.EnableUpdatedAtColumn {
		sqlField := ",\n updated_at " +
			db_schema.TimestampWithTimezone.String() + " " +
			db_schema.Null.String()
		sqlMeta.SqlFieldList = append(sqlMeta.SqlFieldList, sqlField)
	}
	if m.EnableDeleteAtColumn {
		sqlField := ",\n deleted_at " +
			db_schema.TimestampWithTimezone.String() + " " +
			db_schema.Null.String()
		sqlMeta.SqlFieldList = append(sqlMeta.SqlFieldList, sqlField)
	}

	sqlMeta.SqlFieldList = append(sqlMeta.SqlFieldList, "\n);")
}

func (m Migrator) insertFixtures(sqlMeta *dto.Meta) error {
	var sqlFieldList []string
	sqlFieldList = append(sqlFieldList, "INSERT INTO %s \n(")

	var sqlValues []any
	sqlValues = append(sqlValues, sqlMeta.TableName)

	processedCount := 0
	for _, fieldName := range sqlMeta.FieldNames {
		sqlField := "\n %s %s"

		processedCount++
		if processedCount < sqlMeta.FieldCount {
			s += ","
		}

		sqlFieldList = append(sqlFieldList, sqlField)
	}
}
