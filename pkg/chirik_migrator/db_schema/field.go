package db_schema

import (
	"fmt"
	"strings"
)

// TODO: требуется доработка Uuid
const (
	Bool                     DataType = "BOOLEAN"
	Smallint                 DataType = "SMALLINT"
	Int                      DataType = "INTEGER"
	Bigint                   DataType = "BIGINT"
	BigSerial                DataType = "BIGSERIAL"
	Real                     DataType = "REAL"
	DoublePrecision          DataType = "DOUBLE PRECISION"
	Varchar                  DataType = "VARCHAR"
	Text                     DataType = "TEXT"
	TimestampWithTimezone    DataType = "TIMESTAMP WITH TIME ZONE"
	TimestampWithoutTimezone DataType = "TIMESTAMP WITHOUT TIME ZONE"
	Bytes                    DataType = "SMALLINT"
	Uuid                     DataType = "UUID"
	Json                     DataType = "JSON"
	Jsonb                    DataType = "JSONB"
	Null                     DataType = "NULL"
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
	//TODO: это мок, а не полноценная реализация HasIndex
	HasIndex        bool
	Comment         string
	Size            int
	IgnoreMigration bool
	Schema          *Schema
	//TODO: необходимо доработать EmbeddedSchema и OwnerSchema
	EmbeddedSchema *Schema
	OwnerSchema    *Schema
}

type FieldTypeError struct {
	Msg string
}

func (e *FieldTypeError) Error() string {
	return e.Msg
}

// DataType TODO: вынести отдельно
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
	case "bigserial":
		f.DataType = BigSerial
		break
	case "float32", "real":
		f.DataType = Real
		break
	case "float", "float64", "double precision":
		f.DataType = DoublePrecision
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
		break
	case "time", "time.time", "timestamp with time zone":
		f.DataType = TimestampWithTimezone
		break
	case "uuid", "pgtype.uuid":
		f.DataType = Uuid
		break
	case "json":
		f.DataType = Json
		break
	case "jsonb", "pgtype.jsonbcodec":
		f.DataType = Jsonb
		break
	default:
		return f, &FieldTypeError{fmt.Sprintf("unknown type: %s", fieldType)}
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
