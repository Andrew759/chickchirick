package data_type

const (
	typePostgres = "postgres"
	typeMySQL    = "mysql"
)

func PrepareTypeContainer(dbType string) Type {
	switch dbType {
	case typePostgres:
		return PgType("")
	case typeMySQL:
		break
	default:
		return PgType("")
	}
	return PgType("")
}
