package provider

import (
	"chickChirick/pkg/chirik_ast"
	"chickChirick/pkg/chirik_migrator/console/config"
	"chickChirick/pkg/chirik_migrator/db_schema"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type MigratorInfo struct {
	MigratorEnabled bool
	Schema          db_schema.Schema
	ErrList         []error
}

func (mInfo *MigratorInfo) FillByEntity(structure chirik_ast.Structure) {
	//TODO: тут подумать над однообразной обработкой ошибок
	fields := structure.Fields()

	//Сначала отдельно проверяется главный тег
	mainMigratorTag := fields.Tag(config.MigratorTag)
	if mainMigratorTag == nil {
		mInfo.ErrList = append(mInfo.ErrList, fmt.Errorf("main migrator tag %s was not found", config.MigratorTag))
		return
	}

	schema := mInfo.PrepareEmptySchema(structure)
	var schemaFields []*db_schema.Field

	for _, field := range fields.List() {
		schemaField, err := mInfo.prepareSchemaField(*field, &schema)
		if err != nil {
			mInfo.ErrList = append(mInfo.ErrList, err)
		}

		//Пропуск незначащих полей: могут иметь побочные действия, но при непосредственной
		// миграции использоваться не могут
		if schemaField.Name == "" {
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
		}
	}

	mInfo.Schema = schema
}

func (mInfo *MigratorInfo) HasError() bool {
	return mInfo.ErrList != nil
}

func (mInfo *MigratorInfo) PrepareEmptySchema(structure chirik_ast.Structure) db_schema.Schema {
	schema := db_schema.Schema{}
	schema.Name = structure.Name()

	tableName := strings.ToLower(structure.Name())
	schema.Table = tableName

	return schema
}

func (mInfo *MigratorInfo) prepareSchemaField(field chirik_ast.Field, schema *db_schema.Schema) (db_schema.Field, error) {
	schemaField := db_schema.Field{}

	schemaField.Name = field.Name()
	schemaField.Schema = schema

	schemaField, err := schemaField.FillPgDataTypeByString(field.Type().Value())

	fTags := field.Tags()
	if fTags != nil {
		//fTags.ListMock()
		for _, tag := range fTags.List() {
			err := mInfo.fillByTag(tag, &schemaField)
			if err != nil {
				break
			}
		}
	}

	return schemaField, err
}

func (mInfo *MigratorInfo) fillByTag(tag chirik_ast.Tag, schemaField *db_schema.Field) error {
	tKey := tag.Key
	tScalarVal := tag.Values[0]

	switch tKey {
	case config.MigratorTag:
		switch tScalarVal {
		case config.MigratorEnabled:
			mInfo.MigratorEnabled = true
		case config.MigratorDisabled:
		default:
			mInfo.MigratorEnabled = false
		}
	case config.MigratorGormTag:
		return mInfo.fillByGormTag(schemaField, tag.Values)
	}

	return nil
}

func (mInfo *MigratorInfo) fillByGormTag(schemaField *db_schema.Field, tValues []string) error {
	var err error
	tValueWithSizeRe := regexp.MustCompile(`([a-zA-Zа-яА-ЯёЁ]+)\((\d+)\)`)

	for _, tFullValue := range tValues {
		splitTValue := strings.Split(tFullValue, ":")

		tPrefix := strings.Trim(splitTValue[0], `"`)
		tValue := strings.Trim(splitTValue[0], `"`)
		if len(splitTValue) > 1 {
			tValue = strings.Trim(splitTValue[1], `"`)
		}

		valuesWithSize := tValueWithSizeRe.FindStringSubmatch(tValue)
		if len(valuesWithSize) == 3 {
			tValue = valuesWithSize[1]

			//В данном кейсе размер устанавливается заранее
			schemaField.Size, err = strconv.Atoi(valuesWithSize[2])
		}

		switch tPrefix {
		case "column":
			schemaField.Name = tValue
		case "type":
			_, err = schemaField.FillPgDataTypeByString(tValue)
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
		}

		if err != nil {
			break
		}
	}

	return err
}
