package abstraction

import (
	"chickChirick/pkg/chirik_migrator/db_schema"
	migratorDto "chickChirick/pkg/chirik_migrator/migrator/provider"
	"reflect"
)

type Migrator interface {
	// Tables
	CreateTables(entities map[string][]migratorDto.MigratorInfo) error
	CreateTable(entity migratorDto.MigratorInfo) error
	DropTable(entity migratorDto.MigratorInfo) error

	// Columns
	AddColumn(entity migratorDto.MigratorInfo, field string) error
	DropColumn(entity migratorDto.MigratorInfo, field string) error
	AlterColumn(entity migratorDto.MigratorInfo, field string) error
	MigrateColumn(dst interface{}, field *db_schema.Field, columnType ColumnType) error

	// Constraints
	CreateConstraint(dst interface{}, name string) error
	DropConstraint(dst interface{}, name string) error
	HasConstraint(dst interface{}, name string) bool

	// Indexes
	CreateIndex(dst interface{}, name string) error
	DropIndex(dst interface{}, name string) error
}

type ColumnType interface {
	Name() string
	DatabaseTypeName() string                 // varchar
	ColumnType() (columnType string, ok bool) // varchar(64)
	PrimaryKey() (isPrimaryKey bool, ok bool)
	AutoIncrement() (isAutoIncrement bool, ok bool)
	Length() (length int64, ok bool)
	DecimalSize() (precision int64, scale int64, ok bool)
	Nullable() (nullable bool, ok bool)
	Unique() (unique bool, ok bool)
	ScanType() reflect.Type
	Comment() (value string, ok bool)
	DefaultValue() (value string, ok bool)
}
