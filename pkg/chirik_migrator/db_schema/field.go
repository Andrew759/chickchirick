package db_schema

import (
	"fmt"
	"strings"
)

type (
	DataType string
)

const (
	Bool     DataType = "bool"
	Smallint DataType = "smallint"
	Int      DataType = "integer"
	Bigint   DataType = "bigint"
	Float    DataType = "float"
	Varchar  DataType = "varchar"
	Time     DataType = "time"
	Bytes    DataType = "bytes"
	Uuid     DataType = "uuid"
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
	HasIndex       bool
	Comment        string
	Size           int
	Schema         *Schema
	EmbeddedSchema *Schema
	OwnerSchema    *Schema
}

func (f Field) FillPgDataTypeByString(fieldType string) (Field, error) {
	fieldType = strings.ToLower(fieldType)

	switch fieldType {
	case "int8", "int16", "uint8", "uint16":
		f.DataType = Smallint
		break
	case "int32", "uint32":
		f.DataType = Int
		break
	case "int", "int64", "bigint", "uint", "uint64":
		f.DataType = Bigint
		break
	case "float", "float32", "float64":
		f.DataType = Float
		break
	case "string", "varchar":
		f.DataType = Varchar
		break
	case "bool":
		f.DataType = Bool
		break
	case "byte", "rune":
		f.DataType = Bytes
		break
	case "time":
		f.DataType = Time
		break
	case "uuid":
		f.DataType = Uuid
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
