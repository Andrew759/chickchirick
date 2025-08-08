package chirik_faker

import (
	"chickChirick/pkg/chirik_migrator/db_schema/data_type"
	"encoding/json"
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"time"
)

func FakeSQLValue(fieldName string, dataType data_type.Type) (any, error) {
	switch fieldName {
	case "id":
		return gofakeit.IntRange(1, 32767), nil
	case "name":
		return gofakeit.Name(), nil
	case "phone":
		return gofakeit.Phone(), nil
	case "email":
		return gofakeit.Email(), nil
	case "password":
		return gofakeit.Password(true, true, true, true, false, 20), nil
	case "created_at", "updated_at", "deleted_at":
		return gofakeit.Date().Format(time.RFC3339), nil
	case "token":
		return gofakeit.UUID(), nil
	}

	switch dataType {
	case dataType.Bool():
		return gofakeit.Bool(), nil

	case dataType.SmallInt():
		return gofakeit.IntRange(1, 32767), nil

	case dataType.Int():
		return gofakeit.IntRange(1, 2147483647), nil

	case dataType.BigInt():
		return gofakeit.Int64(), nil

	case dataType.BigSerial():
		return gofakeit.IntRange(1, 9223372036854775807), nil

	case dataType.Real():
		return gofakeit.Float32Range(1, 1000), nil

	case dataType.DoublePrecision():
		return gofakeit.Float64Range(1, 1e6), nil

	case dataType.Varchar(), dataType.Text():
		return gofakeit.Sentence(5), nil

	case dataType.TimestampWithTimezone():
		return gofakeit.Date().Format(time.RFC3339), nil

	case dataType.TimestampWithoutTimezone():
		return gofakeit.Date().Format("2006-01-02 15:04:05"), nil

	case dataType.Uuid():
		return gofakeit.UUID(), nil

	case dataType.Json(), dataType.Jsonb():
		fakeMap := map[string]string{
			"name":  gofakeit.FirstName(),
			"email": gofakeit.Email(),
		}
		b, _ := json.Marshal(fakeMap)

		return fmt.Sprintf("%s", string(b)), nil

	case dataType.Null():
		return nil, nil

	default:
		return "", fmt.Errorf("unsupported type: %s", dataType)
	}
}

func FakeStringWithLength(length uint) string {
	return gofakeit.LetterN(length)
}
