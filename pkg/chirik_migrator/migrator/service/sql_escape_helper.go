package service

import (
	"chickChirick/pkg/chirik_migrator/migrator/dto"
	"fmt"
	"strconv"
	"strings"
)

func AddQuotesToStringValInValueList(sql []string) []string {
	for i, v := range sql {
		sql[i] = AddQuotesToStringVal(v)
	}

	return sql
}

func AddQuotesToStringVal(sql string) string {
	if isNumeric(sql) {
		return sql
	}
	return "'" + sql + "'"
}

func isNumeric(s string) bool {
	if _, err := strconv.Atoi(s); err == nil {
		return true
	}
	if _, err := strconv.ParseFloat(s, 64); err == nil {
		return true
	}
	return false
}

// BuildRawSql Заменяет ? в SQL-запросе на экранированные значения
func BuildRawSql(sql string, sqlMeta dto.Meta) string {
	var stringBuilder strings.Builder

	placeHolderIndex := 0
	for i := 0; i < len(sql); i++ {
		if sql[i] == '?' && placeHolderIndex < sqlMeta.FieldCount {
			fieldMeta := sqlMeta.SqlValues[placeHolderIndex]
			if fieldMeta.IsSafe {
				stringBuilder.WriteString(fieldMeta.Value)
			} else {
				stringBuilder.WriteString(escapeSQLValue(fieldMeta.Value))
				placeHolderIndex++
			}
		} else {
			stringBuilder.WriteByte(sql[i])
		}
	}

	if placeHolderIndex != sqlMeta.FieldCount {
		//TODO: выбрость здсь ошибку, что количество плейсхолдеров не соответствует числу обратанных значений
	}

	return stringBuilder.String()
}

func escapeSQLValue(val interface{}) string {
	switch v := val.(type) {
	case nil:
		return "NULL"
	case string:
		// Заменяем одиночные кавычки на двойные одиночные кавычки
		escaped := strings.ReplaceAll(v, "'", "''")
		return "'" + escaped + "'"
	case bool:
		if v {
			return "TRUE"
		}
		return "FALSE"
	case int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64,
		float32, float64:
		return fmt.Sprintf("%v", v)
	default:
		panic(fmt.Sprintf("unsupported type: %T", val))
	}
}
