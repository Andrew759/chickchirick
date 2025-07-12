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

type Config struct {
	//TOOD: CreateIndexAfterCreateTable сейчас не используется. Проверить необходимость
	CreateIndexAfterCreateTable bool
	mainService.DBDecorator
	MigrationFilesPath   string
	EnableTableNamespace bool
	FixtureCount         int
	FixturePrefix        string
	FixtureNilColumns    map[string]struct{}
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
	err = m.processSQLMeta(sqlMeta)
	if err != nil {
		return err
	}

	if m.FixtureCount > 0 {
		return m.InsertFixtures(sqlMeta)
	}

	return nil
}

func (m Migrator) InsertFixtures(sqlMeta dto.Meta) error {
	var errList []error
	for i := 0; i < m.FixtureCount; i++ {
		sqlMetaIterationCopy := sqlMeta
		err := m.writeFixtureToSqlMeta(&sqlMetaIterationCopy)
		if err != nil {
			errList = append(errList, err)
		}
		err = m.processSQLMeta(sqlMetaIterationCopy)
		if err != nil {
			errList = append(errList, err)
		}
	}
	if len(errList) > 0 {
		//TODO: доработать
		return fmt.Errorf("%s", errList)
	}

	return nil
}

func (m Migrator) processSQLMeta(sqlMeta dto.Meta) error {
	var resultSQL string
	for _, sqlField := range sqlMeta.SqlFieldList {
		resultSQL += sqlField
	}

	var err error
	resultSQL, err = service.BuildRawSql(resultSQL, sqlMeta)
	if err != nil {
		return err
	}

	_, err = m.NativeDB().Exec(resultSQL)
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

func (m Migrator) writeFixtureToSqlMeta(sqlMeta *dto.Meta) error {
	var sqlFieldList []string
	sqlFieldList = append(sqlFieldList, "INSERT INTO ? (")

	var sqlValues []dto.ValueMeta
	sqlValues = append(sqlValues, dto.ValueMeta{
		Value:  sqlMeta.TableName,
		IsSafe: true,
	})

	processedCount := 0
	for _, valueMeta := range sqlMeta.SqlValues {
		if !valueMeta.IsValueStore {
			continue
		}
		valName := fmt.Sprintf("%v", valueMeta.Value)
		//При установке фикстуры устанавливается флаг - не безопасно
		valueMeta.IsSafe = false

		//sqlValues = append(sqlValues, valueMeta)
		sqlFieldList = append(sqlFieldList, fmt.Sprintf("%v", valueMeta.Value))
		processedCount++

		if processedCount < sqlMeta.FieldCount {
			sqlFieldList = append(sqlFieldList, ", ")
		} else if processedCount == sqlMeta.FieldCount {
			sqlFieldList = append(sqlFieldList, ")")
		}

		fixtureValue, err := chirik_faker.FakeValue(valName, valueMeta.Type)
		if err != nil {
			return err
		}

		//Если поле сконфигурировано так, чтобы игнорировать фикстуры - в качестве значения устанавливается nil
		if _, hasValue := m.FixtureNilColumns[valName]; hasValue {
			fixtureValue = nil
		}

		sqlValues = append(sqlValues, dto.ValueMeta{
			Value:  fixtureValue,
			IsSafe: false,
		})

	}

	sqlFieldList = append(sqlFieldList, " VALUES (")
	for i := 1; i <= sqlMeta.FieldCount; i++ {
		if i != sqlMeta.FieldCount {
			sqlFieldList = append(sqlFieldList, "?, ")
		} else {
			sqlFieldList = append(sqlFieldList, "?")
		}
	}
	sqlFieldList = append(sqlFieldList, ");")

	sqlMeta.SqlFieldList = sqlFieldList
	sqlMeta.SqlValues = sqlValues
	sqlMeta.FieldCommentList = []string{}
	sqlMeta.FieldCommentValues = []string{}

	if m.FixturePrefix != "" {
		sqlMeta.MigrationPrefix = m.FixturePrefix + "_" + sqlMeta.MigrationPrefix
	}

	return nil
}

// DropTable TODO: implement this
func (m Migrator) DropTable(entity migratorDto.MigratorInfo) error {
	return nil
}

// CreateConstraint TODO: implement this
func (m Migrator) CreateConstraint(entity migratorDto.MigratorInfo, name string) error {
	return nil
}

// DropConstraint TODO: implement this
func (m Migrator) DropConstraint(entity migratorDto.MigratorInfo, name string) error {
	return nil
}

// CreateIndex TODO: implement this
func (m Migrator) CreateIndex(entity migratorDto.MigratorInfo, name string) error {
	return nil
}

// DropIndex TODO: implement this
func (m Migrator) DropIndex(entity migratorDto.MigratorInfo, name string) error {
	return nil
}
