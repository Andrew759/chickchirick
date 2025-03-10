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
	MigratorOrmType  = "c_migrator_orm"
	MigratorEnabled  = "enabled"
	MigratorDisabled = "disabled"
)

// Список поддерживаемых ORM
const (
	GORM = iota
	BUN
)

const (
	MigrationPath = "DB_MIGRATION_PATH"
	EntityPath    = "DB_ENTITY_PATH"
)
