package configuration

// environment
const (
	Enviroment = "ENVIRONMENT"
	Dev        = "DEV"
	Test       = "TEST"
	Prod       = "PROD"
)

// database
const (
	DbHost           = "DB_HOST"
	DbPort           = "DB_PORT"
	DbName           = "DB_NAME"
	DbUser           = "DB_USER"
	DbPass           = "DB_PASS"
	DbMigrationPatch = "DB_MIGRATION_PATCH"
	DbTimezone       = "DB_TIMEZONE"
)

// redis
const (
	RedisHost     = "REDIS_HOST"
	RedisPort     = "REDIS_PORT"
	RedisUser     = "REDIS_USER"
	RedisPassword = "REDIS_PASSWORD"
)
