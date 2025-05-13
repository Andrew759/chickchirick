package dto

import (
	"chickChirick/pkg/chirik_ast"
	"chickChirick/pkg/chirik_migrator/console/config"
	"chickChirick/pkg/chirik_migrator/db_schema"
	"fmt"
	"strings"
)

type MigratorInfo struct {
	MigratorEnabled bool
	//TODO: дорабоать
	//EntityInfo      dto.FileInfo
	Schema  db_schema.Schema
	ErrList []error
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
	for _, field := range fields.List() {
		schemaField, err := mInfo.prepareSchemaField(*field)
		if err != nil {
			mInfo.ErrList = append(mInfo.ErrList, err)
		}

		//TODO: временно
		fmt.Println(schema, schemaField)
	}

	//TODO: тут предусмотреть удобную структуру с считанными тегами для мигратора. А это можно удалить
	//mInfo.EntityInfo = fileInfo
}

func (mInfo *MigratorInfo) HasError() bool {
	return mInfo.ErrList != nil
}

func (mInfo *MigratorInfo) PrepareEmptySchema(structure chirik_ast.Structure) db_schema.Schema {
	tableName := strings.ToLower(structure.Name())

	schema := db_schema.Schema{}
	schema.Name = tableName

	return schema
}

func (mInfo *MigratorInfo) prepareSchemaField(field chirik_ast.Field) (db_schema.Field, error) {
	schemaField := db_schema.Field{}

	fName := field.Name()
	fType := field.Type().Value()
	fTags := field.Tags()

	schemaField.Name = fName
	schemaField, err := schemaField.FllDataTypeByString(fType)

	if fTags != nil {
		for _, tag := range fTags.List() {
			mInfo.fillByTag(tag, &schemaField)
		}
	}

	return schemaField, err
}

func (mInfo *MigratorInfo) fillByTag(tag chirik_ast.Tag, schemaField *db_schema.Field) {
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
		mInfo.fillByGormTag(*schemaField, tag.Values)
	}
}

func (mInfo *MigratorInfo) fillByGormTag(schemaField db_schema.Field, tValues []string) {
	for _, tValue := range tValues {
		switch tValue {
		case "column":
		case "type":
		case "size":
		case "primaryKey":
		case "unique":
		case "default":
		case "not null":
		case "autoincrement":
		case "autoIncrementIncrement":
		case "index":
		case "uniqueIndex":
		}
	}
}
