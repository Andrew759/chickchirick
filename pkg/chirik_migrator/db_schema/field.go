package db_schema

import (
	"fmt"
	"strings"
)

type (
	DataType string
)

const (
	Bool   DataType = "bool"
	Int    DataType = "int"
	Uint   DataType = "uint"
	Float  DataType = "float"
	String DataType = "string"
	Time   DataType = "time"
	Bytes  DataType = "bytes"
	Uuid   DataType = "uuid"
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

func (f Field) FllDataTypeByString(fieldType string) (Field, error) {
	fieldType = strings.ToLower(fieldType)

	switch fieldType {
	case "int", "int8", "int16", "int32", "int64":
		f.DataType = Int
		break
	case "uint", "uint8", "uint16", "uint32", "uint64":
		f.DataType = Uint
		break
	case "float", "float32", "float64":
		f.DataType = Float
		break
	case "string":
		f.DataType = String
		break
	case "bool":
		f.DataType = Bool
		break
	case "byte", "rune":
		f.DataType = Bytes
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
