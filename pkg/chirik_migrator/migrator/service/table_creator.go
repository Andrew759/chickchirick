package service

import (
	"chickChirick/pkg/chirik_migrator/migrator/dto"
	migratorDto "chickChirick/pkg/chirik_migrator/migrator/provider"
	"strconv"
)

type TableCreator struct {
	SqlProcessor
}

func (tc TableCreator) CreateTables(migratorEntities map[string][]migratorDto.MigratorInfo) ([][]dto.Meta, error) {
	var processedTablesSqlMeta [][]dto.Meta
	for _, migratorInfoList := range migratorEntities {
		for _, migratorInfo := range migratorInfoList {
			processedSqlMetas, err := tc.CreateTable(migratorInfo)
			if err != nil {
				return processedTablesSqlMeta, err
			}
			processedTablesSqlMeta = append(processedTablesSqlMeta, processedSqlMetas)
		}
	}
	return processedTablesSqlMeta, nil
}

func (tc TableCreator) CreateTable(migratorInfo migratorDto.MigratorInfo) ([]dto.Meta, error) {
	var processedSqlMetas []dto.Meta

	err := ValidateMInfo(migratorInfo)
	if err != nil {
		return processedSqlMetas, err
	}

	sqlMeta := tc.processSchemaFields(migratorInfo)
	err = tc.ProcessSQLMeta(sqlMeta)
	if err != nil {
		return processedSqlMetas, err
	}

	processedSqlMetas = append(processedSqlMetas, sqlMeta)

	return processedSqlMetas, nil
}

func (tc TableCreator) processSchemaFields(migratorInfo migratorDto.MigratorInfo) dto.Meta {
	var sqlFieldList []string
	sqlFieldList = append(sqlFieldList, "CREATE TABLE IF NOT EXISTS ? \n(")

	schema := &migratorInfo.Schema
	fullTableName := schema.Table
	if tc.MConfig.EnableTableNamespace {
		fullTableName = migratorInfo.EntityNamespace + "_" + schema.Table
	}

	var sqlValues []dto.ValueMeta
	sqlValues = append(sqlValues, dto.ValueMeta{
		Value:  fullTableName,
		IsSafe: true,
	})

	//TODO: доработать наподобие с sqlValues
	var fieldCommentList []string
	var fieldCommentValues []string
	fieldCount := len(schema.Fields)
	processedCount := 0

	for _, field := range schema.Fields {
		fieldType := field.DataType
		if fieldType == nil {
			continue
		}
		//Пропуск полей без типа
		if !field.HasDataType() {
			continue
		}

		sqlField := "\n ? ?"

		//TODO: работает только для постгры
		if field.AutoIncrement {
			fieldType = field.DataType.BigSerial()
		}
		sqlValues = append(sqlValues,
			dto.ValueMeta{
				Value:        field.Name,
				Type:         fieldType,
				IsSafe:       true,
				IsValueStore: true,
			},
			dto.ValueMeta{
				Value:  fieldType,
				IsSafe: true,
			})

		if field.Size != 0 {
			sqlField += "(?)"
			sqlValues = append(sqlValues,
				dto.ValueMeta{
					Value:  strconv.Itoa(field.Size),
					Type:   fieldType,
					IsSafe: true,
				})
		}

		if field.PrimaryKey {
			sqlField += " PRIMARY KEY"
		}
		if field.HasDefaultValue {
			sqlField += " DEFAULT ?"
			sqlValues = append(sqlValues,
				dto.ValueMeta{
					Value:  field.DefaultValue,
					Type:   fieldType,
					IsSafe: true,
				})

		}
		if field.NotNull {
			sqlField += " NOT NULL"
		}
		if field.Unique {
			sqlField += " UNIQUE"
		}
		if field.Comment != "" {
			fieldComment := "comment on column ?.? is '?';"
			fieldCommentValues = append(fieldCommentValues, schema.Table, field.Name, field.Comment)
			fieldCommentList = append(fieldCommentList, fieldComment)
		}

		processedCount++
		if processedCount < fieldCount {
			sqlField += ","
		}

		sqlFieldList = append(sqlFieldList, sqlField)
	}

	sqlFieldList = append(sqlFieldList, "\n);")

	return dto.Meta{
		TableName:          fullTableName,
		SqlFieldList:       sqlFieldList,
		SqlValues:          sqlValues,
		FieldCommentList:   fieldCommentList,
		FieldCommentValues: fieldCommentValues,
		FieldCount:         fieldCount,
		MigrationPrefix:    fullTableName,
	}
}
