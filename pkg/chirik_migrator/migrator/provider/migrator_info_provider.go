package provider

import (
	"chickChirick/pkg/chirik_ast"
	"chickChirick/pkg/chirik_migrator/console/config"
	"chickChirick/pkg/chirik_migrator/db_schema"
	"chickChirick/pkg/chirik_migrator/migrator/service"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type MigratorInfo struct {
	MigratorEnabled bool
	EntityNamespace string
	Schema          db_schema.Schema
	ErrList         []error
	InfoErrList     []error
}

func (mInfo *MigratorInfo) FillByEntity(structure chirik_ast.Structure) {
	fields := structure.Fields()

	//Сначала отдельно проверяется главный тег
	mainMigratorTag := fields.Tag(config.MigratorTag)
	if mainMigratorTag == nil {
		mInfo.ErrList = append(mInfo.ErrList, fmt.Errorf(
			"main migrator tag was not found: %s", config.MigratorTag),
		)
		return
	}

	schema := mInfo.PrepareEmptySchema(structure)
	var schemaFields []*db_schema.Field

	for _, field := range fields.List() {
		schemaField := mInfo.PrepareSchemaField(*field, &schema)
		//Пропуск незначащих полей: могут иметь побочные действия, но при непосредственной
		// миграции использоваться не могут
		if schemaField.Name == "" ||
			schemaField.DataType.String() == "" ||
			schemaField.IgnoreMigration {
			continue
		}

		schemaFields = append(schemaFields, &schemaField)
	}

	schema.Fields = schemaFields

	//TODO: сейчас не работает
	var schemaForeignKey []*db_schema.Field
	schema.ForeignKey = schemaForeignKey

	for _, preparedField := range schemaFields {
		if preparedField.PrimaryKey {
			schema.PrimaryKey = preparedField
			break
		}
	}

	mInfo.Schema = schema
}

func (mInfo *MigratorInfo) HasCriticalError() bool {
	return mInfo.ErrList != nil
}

func (mInfo *MigratorInfo) HasInfoError() bool {
	return mInfo.InfoErrList != nil
}

func (mInfo *MigratorInfo) PrepareEmptySchema(structure chirik_ast.Structure) db_schema.Schema {
	schema := db_schema.Schema{}
	schema.Name = structure.Name()

	tableName := service.AddSingleSPostfix(
		service.ToSnakeCase(structure.Name()),
	)

	schema.Table = tableName

	return schema
}

func (mInfo *MigratorInfo) PrepareSchemaField(field chirik_ast.Field, schema *db_schema.Schema) db_schema.Field {
	schemaField := db_schema.Field{}

	schemaField.Name = service.ToSnakeCase(field.Name())
	schemaField.Schema = schema

	var infoErr error
	schemaField, infoErr = schemaField.FillPgDataTypeByString(field.Type().Value())
	if infoErr != nil {
		mInfo.InfoErrList = append(mInfo.ErrList, infoErr)
	}

	fTags := field.Tags()
	if fTags != nil {
		for _, tag := range fTags.List() {
			mInfo.FillByTagAndSchemaField(tag, &schemaField)
			if mInfo.HasCriticalError() {
				break
			}
		}
	}

	return schemaField
}

func (mInfo *MigratorInfo) FillByTagAndSchemaField(tag chirik_ast.Tag, schemaField *db_schema.Field) {
	tKey := tag.Key
	tScalarVal := tag.Values[0]

	switch tKey {
	case config.MigratorTag:
		switch tScalarVal {
		case config.MigratorEnabled:
			mInfo.MigratorEnabled = true
		case config.MigratorDisabled:
			mInfo.MigratorEnabled = false
		default:
			mInfo.MigratorEnabled = false
		}
	case config.MigratorGormTag:
		mInfo.FillByGormTagAndSchemaField(schemaField, tag.Values)
	}
}

func (mInfo *MigratorInfo) FillByGormTagAndSchemaField(schemaField *db_schema.Field, tValues []string) {
	tValueWithSizeRe := regexp.MustCompile(`([a-zA-Zа-яА-ЯёЁ]+)\((\d+)\)`)
	var err error

	for _, tFullValue := range tValues {
		splitTValue := strings.Split(tFullValue, ":")
		tPrefix := strings.Trim(splitTValue[0], `"`)

		//Если значение не составное, то по умолчанию значением будет являться префикс
		tValue := strings.Trim(splitTValue[0], `"`)
		if len(splitTValue) > 1 {
			tValue = strings.Trim(splitTValue[1], `"`)
		}

		valuesWithSize := tValueWithSizeRe.FindStringSubmatch(tValue)
		if len(valuesWithSize) == 3 {
			tValue = valuesWithSize[1]

			//В данном кейсе размер устанавливается заранее
			schemaField.Size, err = strconv.Atoi(valuesWithSize[2])
			if err != nil {
				mInfo.ErrList = append(mInfo.ErrList, err)
			}
		}

		switch tPrefix {
		case "column":
			schemaField.Name = tValue
		case "type":
			_, typeErr := schemaField.FillPgDataTypeByString(tValue)
			mInfo.InfoErrList = append(mInfo.InfoErrList, typeErr)
		case "size":
			schemaField.Size, err = strconv.Atoi(tValue)
		case "primaryKey":
			schemaField.PrimaryKey = true
		case "unique":
			schemaField.Unique = true
		case "default":
			schemaField.DefaultValue = tValue
		case "not null":
			schemaField.NotNull = true
		case "autoIncrement":
			schemaField.AutoIncrement = true
		case "autoIncrementIncrement":
			schemaField.AutoIncrementIncrement, err = strconv.ParseInt(tValue, 10, 64)
		case "index":
			//TODO: не реализовано
			schemaField.HasIndex = true
		case "uniqueIndex":
			schemaField.Unique = true
		case "comment":
			schemaField.Comment = tValue
		case "ignoreMigration":
			schemaField.IgnoreMigration = true
		}

		if err != nil {
			mInfo.ErrList = append(mInfo.ErrList, err)
			break
		}
	}
}
