package db_schema

import (
	"chickChirick/pkg/chirik_migrator/db_schema/data_type"
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

func (f Field) FillDataTypeByString(fieldType string) (Field, error) {
	fieldType = strings.ToLower(fieldType)

	switch fieldType {
	case "smallint", "int8", "int16", "uint8", "uint16":
		f.DataType = f.DataType.SmallInt()
		break
	case "int32", "uint32", "int", "uint":
		f.DataType = f.DataType.Int()
		break
	case "int64", "uint64", "bigint":
		f.DataType = f.DataType.BigInt()
		break
	case "bigserial":
		f.DataType = f.DataType.BigSerial()
		break
	case "float32", "real":
		f.DataType = f.DataType.Real()
		break
	case "float", "float64", "double precision":
		f.DataType = f.DataType.DoublePrecision()
		break
	case "string", "varchar":
		f.DataType = f.DataType.Varchar()
		break
	case "text":
		f.DataType = f.DataType.Text()
	case "bool", "boolean":
		f.DataType = f.DataType.Bool()
		break
	case "byte", "rune":
		f.DataType = f.DataType.Bytes()
		break
	case "timestamp without time zone":
		f.DataType = f.DataType.TimestampWithoutTimezone()
		break
	case "time", "time.time", "timestamp with time zone":
		f.DataType = f.DataType.TimestampWithTimezone()
		break
	case "uuid", "pgtype.uuid":
		f.DataType = f.DataType.Uuid()
		break
	case "json":
		f.DataType = f.DataType.Json()
		break
	case "jsonb", "pgtype.jsonbcodec":
		f.DataType = f.DataType.Jsonb()
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
	return f.DataType.IsHasSetValue()
}

func (f Field) HasSize() bool {
	return f.Size != 0
}
