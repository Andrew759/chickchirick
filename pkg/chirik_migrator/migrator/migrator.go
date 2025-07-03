package migrator

import (
	mainService "chickChirick/cmd/service"
	"chickChirick/pkg/chirik_faker"
	"chickChirick/pkg/chirik_migrator/db_schema"
	"chickChirick/pkg/chirik_migrator/file"
	"chickChirick/pkg/chirik_migrator/migrator/dto"
	migratorDto "chickChirick/pkg/chirik_migrator/migrator/provider"
	"chickChirick/pkg/chirik_migrator/migrator/service"
	"errors"
	"fmt"
	"strconv"
)

// TODO: согласовать с интерфейсом abstraction/Migrator
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
	FixturePrefix         string
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

	err = m.processSQLMeta(sqlMeta)
	if err != nil {
		return err
	}

	if m.EnableFixtures {
		err = m.writeFixtureToSqlMeta(&sqlMeta)
		if err != nil {
			return err
		}
		return m.processSQLMeta(sqlMeta)
	}

	return nil
}

func (m Migrator) processSQLMeta(sqlMeta dto.Meta) error {
	var resultSQL string
	for _, sqlField := range sqlMeta.SqlFieldList {
		resultSQL += sqlField
	}

	resultSQL = service.BuildRawSql(resultSQL, sqlMeta)

	_, err := m.NativeDB().Exec(resultSQL)
	if err != nil {
		return err
	}

	return file.WriteSQLToFile(resultSQL, sqlMeta.MigrationPrefix, m.MigrationFilesPath)
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

// TODO: реализовать отдельные провайдеры для разных БД. Сейчас работает только Postgres
func (m Migrator) processSchemaFields(migratorInfo migratorDto.MigratorInfo) dto.Meta {
	var sqlFieldList []string
	sqlFieldList = append(sqlFieldList, "CREATE TABLE IF NOT EXISTS ? \n(")

	schema := &migratorInfo.Schema
	fullTableName := schema.Table
	if m.Config.EnableTableNamespace {
		fullTableName = migratorInfo.EntityNamespace + "_" + schema.Table
	}

	var sqlValues []dto.ValueMeta
	sqlValues = append(sqlValues, dto.ValueMeta{
		Value:  fullTableName,
		Type:   db_schema.Varchar.String(),
		IsSafe: true,
	})

	//TODO: доработать наподобие с sqlValues
	var fieldCommentList []string
	var fieldCommentValues []string
	fieldCount := len(schema.Fields)
	processedCount := 0

	for _, field := range schema.Fields {
		fieldType := field.DataType.String()
		//Пропуск полей без типа
		if fieldType == "" {
			continue
		}

		sqlField := "\n ? ?"

		if field.AutoIncrement {
			fieldType = string(db_schema.BigSerial)
		}
		sqlValues = append(sqlValues,
			dto.ValueMeta{
				Value:  field.Name,
				Type:   db_schema.Varchar.String(),
				IsSafe: true,
			},
			dto.ValueMeta{
				Value:  fieldType,
				Type:   db_schema.Varchar.String(),
				IsSafe: true,
			})

		if field.Size != 0 {
			sqlField += "(?)"
			sqlValues = append(sqlValues,
				dto.ValueMeta{
					Value:  strconv.Itoa(field.Size),
					Type:   db_schema.Int.String(),
					IsSafe: true,
				})
		}

		if field.PrimaryKey {
			sqlField += " PRIMARY KEY"
		}
		if field.HasDefaultValue {
			sqlField += "DEFAULT ?"
			sqlValues = append(sqlValues,
				dto.ValueMeta{
					Value:  field.DefaultValue,
					Type:   db_schema.Varchar.String(),
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
			fieldCommentValues = append(fieldCommentValues, schema.Name, field.Name, field.Comment)
			fieldCommentList = append(fieldCommentList, fieldComment)
		}

		processedCount++
		if processedCount < fieldCount {
			sqlField += ","
		}

		sqlFieldList = append(sqlFieldList, sqlField)
	}

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

// TODO: не попадают в фикстуры, доработать
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

// TODO: сломано
func (m Migrator) writeFixtureToSqlMeta(sqlMeta *dto.Meta) error {
	var sqlFieldList []string
	sqlFieldList = append(sqlFieldList, "INSERT INTO ? (")

	var sqlValues []dto.ValueMeta
	sqlValues = append(sqlValues, dto.ValueMeta{
		Value:  sqlMeta.TableName,
		Type:   db_schema.Varchar.String(),
		IsSafe: true,
	})
	processedCount := 0

	for i, valueMeta := range sqlMeta.SqlValues {
		//Фикстурам требуется экранирование
		valueMeta.IsSafe = false
		sqlMeta.SqlValues[i] = valueMeta

		sqlField := valueMeta.Value
		processedCount++
		if processedCount < sqlMeta.FieldCount {
			sqlField += ", "
		} else if processedCount == sqlMeta.FieldCount {
			sqlField += ")"
		}

		fixtureValue, err := chirik_faker.FakeValue(valueMeta.Value, valueMeta.Type)
		if err != nil {
			return err
		}

		sqlValues = append(sqlValues, dto.ValueMeta{
			Value:  fixtureValue,
			Type:   db_schema.Varchar.String(),
			IsSafe: false,
		})

		sqlFieldList = append(sqlFieldList, sqlField)
	}

	fixtureValueHolder := " VALUES ("
	for i := 1; i <= sqlMeta.FieldCount; i++ {
		if i != sqlMeta.FieldCount {
			fixtureValueHolder += "?, "
		} else {
			fixtureValueHolder += "?"
		}
	}
	fixtureValueHolder += ");"

	sqlFieldList = append(sqlFieldList, fixtureValueHolder)

	sqlMeta.SqlFieldList = sqlFieldList
	sqlMeta.SqlValues = sqlValues
	sqlMeta.FieldCommentList = []string{}
	sqlMeta.FieldCommentValues = []string{}

	if m.FixturePrefix != "" {
		sqlMeta.MigrationPrefix = m.FixturePrefix + "_" + sqlMeta.MigrationPrefix
	}

	return nil
}
