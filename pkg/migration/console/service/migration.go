package service

import "fmt"

type Attribute struct {
	Name         string `json:"name"`
	NameDb       string `json:"name_db"`
	Type         string `json:"type"`
	Default      bool   `json:"default"`
	DefaultValue string `json:"default_value"`
	Length       string `json:"length"`
	Nullable     bool   `json:"nullable"`
	Pk           bool   `json:"pk"`
	Fk           bool   `json:"fk"`
	FkTable      string `json:"fk_table"`
	FkColumn     string `json:"fk_column"`
}

func buildSqlByAttribute(a Attribute) string {
	typeMap := map[string]string{
		"string":   fmt.Sprintf("varchar(%s)", a.Length),
		"text":     "text",
		"json":     "jsonb",
		"uuid":     "uuid",
		"datetime": "timestamp(0)",
		"int":      "int",
		"bool":     "bool",
	}

	sqlType, ok := typeMap[a.Type]
	if !ok {
		panic(fmt.Sprintf("unknown sql type %s", a.Type))
	}

	var sqlDefault string
	if a.Default {
		if a.Type == "int" && a.Pk {
			sqlType = "serial"
		} else if a.DefaultValue == "" {
			sqlDefault = "DEFAULT NULL"
		} else {
			sqlDefault = fmt.Sprintf("DEFAULT %s", a.DefaultValue)
		}
	}

	sql := fmt.Sprintf("%s %s", a.NameDb, sqlType)

	if sqlDefault != "" {
		sql = fmt.Sprintf("%s %s", sql, sqlDefault)
	}

	if a.Pk {
		sql = fmt.Sprintf("%s PRIMARY KEY", sql)
	} else if !a.Nullable {
		sql = fmt.Sprintf("%s NOT NULL", sql)
	}

	if a.Fk {
		sql += fmt.Sprintf(" REFERENCES %s (%s)", a.FkTable, a.FkColumn)
	}

	return sql
}
