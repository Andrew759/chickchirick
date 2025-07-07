package config

// Ключи мигратора
const (
	DBType                = "DB_TYPE"
	MigrationPath         = "DB_MIGRATION_PATH"
	EntityPath            = "DB_ENTITY_PATH"
	EnableTableNamespace  = "ENABLE_TABLE_NAMESPACE"
	EnableCreatedAtColumn = "ENABLE_CREATED_AT_COLUMN"
	EnableUpdatedAtColumn = "ENABLE_CREATED_AT_COLUMN"
	EnableDeleteAtColumn  = "ENABLE_DELETE_AT_COLUMN"
	EnableFixtures        = "ENABLE_FIXTURES"
	FixturePrefix         = "FIXTURE_PREFIX"
)

// Ключи тегов мигратора
const (
	MigratorTag       = "c_migrator"
	MigratorEnabled   = "enabled"
	MigratorDisabled  = "disabled"
	MigratorGormTag   = "gorm"
	MigratorTableName = "c_migrator_t_name"
)
