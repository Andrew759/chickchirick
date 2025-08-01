package helper

import (
	"chickChirick/pkg/chirik_migrator/migrator/dto"
	"fmt"
	"strings"
)

// BuildRawSql Заменяет ? в SQL-запросе на экранированные значения
func BuildRawSql(sql string, sqlMeta dto.Meta) (string, error) {
	var stringBuilder strings.Builder

	placeHolderIndex := 0
	for i := 0; i < len(sql); i++ {
		if sql[i] == '?' && placeHolderIndex < len(sqlMeta.SqlValues) {
			valueMeta := sqlMeta.SqlValues[placeHolderIndex]
			vMetaString := fmt.Sprintf("%v", valueMeta.Value)

			if valueMeta.IsSafe {
				stringBuilder.WriteString(vMetaString)
			} else {
				escapedVal, err := escapeSQLValue(valueMeta.Value)
				if err != nil {
					return "", err
				}
				stringBuilder.WriteString(escapedVal)
			}
			placeHolderIndex++
		} else {
			stringBuilder.WriteByte(sql[i])
		}
	}

	if placeHolderIndex != sqlMeta.FieldCount {
		//TODO: выбросить здесь ошибку, что количество плейсхолдеров не соответствует числу обработанных значений
	}

	return stringBuilder.String(), nil
}

func escapeSQLValue(val interface{}) (string, error) {
	switch v := val.(type) {
	case nil:
		return "NULL", nil
	case string:
		escaped := strings.ReplaceAll(v, "'", "''")
		return "'" + escaped + "'", nil
	case bool:
		if v {
			return "TRUE", nil
		}
		return "FALSE", nil
	case int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64,
		float32, float64:
		return fmt.Sprintf("%v", v), nil
	default:
		return "", fmt.Errorf("unsupported type: %T", val)
	}
}
