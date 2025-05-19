package db_schema

import (
	"fmt"
	"strings"
)

//TODO: впоследствии можно вынести эти константы и методы для posgres в отдельное место.

const (
	Bool                     DataType = "BOOLEAN"
	Smallint                 DataType = "SMALLINT"
	Int                      DataType = "INTEGER"
	Bigint                   DataType = "BIGINT"
	Float                    DataType = "FLOAT"
	Varchar                  DataType = "VARCHAR"
	Text                     DataType = "TEXT"
	TimestampWithTimezone    DataType = "TIMESTAMP WITH TIME ZONE"
	TimestampWithoutTimezone DataType = "TIMESTAMP WITHOUT TIME ZONE"
	Bytes                    DataType = "SMALLINT"
	// Uuid TODO: требуется доработка:
	Uuid  DataType = "UUID"
	Json  DataType = "JSON"
	Jsonb DataType = "JSONB"
)

type Field struct {
	Name                   string
	DataType               DataType
	PrimaryKey             bool
	AutoIncrement          bool
	AutoIncrementIncrement int64
	HasDefaultValue        bool
	DefaultValue           string
	NotNull                bool
	Unique                 bool
	//TODO: это мок, а не полноценная реализация
	HasIndex        bool
	Comment         string
	Size            int
	IgnoreMigration bool
	Schema          *Schema
	//TODO: необходимо доработать
	EmbeddedSchema *Schema
	OwnerSchema    *Schema
}

type DataType string

func (d DataType) String() string {
	return string(d)
}

func (f Field) FillPgDataTypeByString(fieldType string) (Field, error) {
	fieldType = strings.ToLower(fieldType)

	switch fieldType {
	case "smallint", "int8", "int16", "uint8", "uint16":
		f.DataType = Smallint
		break
	case "int32", "uint32":
		f.DataType = Int
		break
	case "int", "int64", "bigint", "uint", "uint64":
		f.DataType = Bigint
		break
		//TODO: не доработано
	case "float", "float32", "float64":
		f.DataType = Float
		break
	case "string", "varchar":
		f.DataType = Varchar
		break
	case "text":
		f.DataType = Text
	case "bool", "boolean":
		f.DataType = Bool
		break
	case "byte", "rune":
		f.DataType = Bytes
		break
	case "timestamp without time zone":
		f.DataType = TimestampWithoutTimezone
	case "time", "time.time", "timestamp with time zone":
		f.DataType = TimestampWithTimezone
		break
	case "uuid", "pgtype.uuid":
		f.DataType = Uuid
	case "json":
		f.DataType = Json
	case "jsonb", "pgtype.jsonbcodec":
		f.DataType = Jsonb
		break
	default:
		return f, fmt.Errorf("unknown type: %s", fieldType)
	}

	return f, nil
}

func (f Field) HasName() bool {
	return f.Name != ""
}

func (f Field) HasDataType() bool {
	return f.DataType != ""
}

func (f Field) HasSize() bool {
	return f.Size != 0
}
