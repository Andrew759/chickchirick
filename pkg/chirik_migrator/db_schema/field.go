package db_schema

import (
	"chickChirick/pkg/chirik_migrator/db_schema/data_type"
	"chickChirick/pkg/chirik_migrator/db_schema/data_type/postgres"
	"fmt"
	"strings"
)

type Field struct {
	Name                   string
	DataType               data_type.Type
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

func (f Field) FillPgDataTypeByString(fieldType string) (Field, error) {
	fieldType = strings.ToLower(fieldType)

	switch fieldType {
	case "smallint", "int8", "int16", "uint8", "uint16":
		f.DataType = postgres.Smallint
		break
	case "int32", "uint32":
		f.DataType = postgres.Int
		break
	case "int", "int64", "bigint", "uint", "uint64":
		f.DataType = postgres.Bigint
		break
	case "bigserial":
		f.DataType = postgres.BigSerial
		break
	case "float32", "real":
		f.DataType = postgres.Real
		break
	case "float", "float64", "double precision":
		f.DataType = postgres.DoublePrecision
		break
	case "string", "varchar":
		f.DataType = postgres.Varchar
		break
	case "text":
		f.DataType = postgres.Text
	case "bool", "boolean":
		f.DataType = postgres.Bool
		break
	case "byte", "rune":
		f.DataType = postgres.Bytes
		break
	case "timestamp without time zone":
		f.DataType = postgres.TimestampWithoutTimezone
		break
	case "time", "time.time", "timestamp with time zone":
		f.DataType = postgres.TimestampWithTimezone
		break
	case "uuid", "pgtype.uuid":
		f.DataType = postgres.Uuid
		break
	case "json":
		f.DataType = postgres.Json
		break
	case "jsonb", "pgtype.jsonbcodec":
		f.DataType = postgres.Jsonb
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
	return f.DataType != nil
}

func (f Field) HasSize() bool {
	return f.Size != 0
}
