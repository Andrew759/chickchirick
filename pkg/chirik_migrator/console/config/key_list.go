package config

// Ключи консольных команд
const (
	MigrateKey     = "migrate"
	AssignmentKey  = "="
	ValueSeparator = " "
	AllTables      = "*"
)

// Ключи тегов мигратора
const (
	MigratorTag      = "c_migrator"
	MigratorEnabled  = "enabled"
	MigratorDisabled = "disabled"
	MigratorGormTag  = "gorm"
)

const (
	MigrationPath = "DB_MIGRATION_PATH"
	EntityPath    = "DB_ENTITY_PATH"
)
