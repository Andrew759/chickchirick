package config

// Ключи консольных команд
const (
	MigrateKey     = "migrate"
	AssignmentKey  = "="
	ValueSeparator = " "
	AllTables      = "*"
)

// Ключи мигратора
const (
	MigrationPath         = "DB_MIGRATION_PATH"
	EntityPath            = "DB_ENTITY_PATH"
	EnableTableNamespace  = "ENABLE_TABLE_NAMESPACE"
	EnableDeleteAtColumn  = "ENABLE_DELETE_AT_COLUMN"
	EnableCreatedAtColumn = "ENABLE_CREATED_AT_COLUMN"
)

// Ключи тегов мигратора
const (
	MigratorTag      = "c_migrator"
	MigratorEnabled  = "enabled"
	MigratorDisabled = "disabled"
	MigratorGormTag  = "gorm"
)
