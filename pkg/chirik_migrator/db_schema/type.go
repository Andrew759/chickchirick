package db_schema

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
	Comment                string
	Size                   int
	Schema                 *Schema
	EmbeddedSchema         *Schema
	OwnerSchema            *Schema
}
