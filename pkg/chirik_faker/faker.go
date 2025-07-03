package chirik_faker

import (
	dbs "chickChirick/pkg/chirik_migrator/db_schema"
	"encoding/json"
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"time"
)

func FakeValue(fieldName string, dataType string) (string, error) {
	switch fieldName {
	case "id":
		return fmt.Sprintf("%d", gofakeit.IntRange(1, 32767)), nil
	case "name":
		return gofakeit.Name(), nil
	case "phone":
		return gofakeit.Phone(), nil
	case "email":
		return gofakeit.Email(), nil
	case "password":
		return gofakeit.Password(true, true, true, true, false, 20), nil
	case "created_at", "updated_at", "deleted_at":
		return gofakeit.TimeZoneFull(), nil
	case "token":
		return gofakeit.UUID(), nil
	}

	switch dataType {
	case dbs.Bool.String():
		return fmt.Sprintf("%v", gofakeit.Bool()), nil

	case dbs.Smallint.String():
		return fmt.Sprintf("%d", gofakeit.IntRange(1, 32767)), nil

	case dbs.Int.String():
		return fmt.Sprintf("%d", gofakeit.IntRange(1, 2147483647)), nil

	case dbs.Bigint.String():
		return fmt.Sprintf("%d", gofakeit.Int64()), nil

	case dbs.BigSerial.String():
		return fmt.Sprintf("%d", gofakeit.IntRange(1, 9223372036854775807)), nil

	case dbs.Real.String():
		return fmt.Sprintf("%.4f", gofakeit.Float32Range(1, 1000)), nil

	case dbs.DoublePrecision.String():
		return fmt.Sprintf("%.6f", gofakeit.Float64Range(1, 1e6)), nil

	case dbs.Varchar.String(), dbs.Text.String():
		return fmt.Sprintf("'%s'", gofakeit.Sentence(5)), nil

	case dbs.TimestampWithTimezone.String():
		return fmt.Sprintf("'%s'", gofakeit.Date().Format(time.RFC3339)), nil

	case dbs.TimestampWithoutTimezone.String():
		return fmt.Sprintf("'%s'", gofakeit.Date().Format("2006-01-02 15:04:05")), nil

	case dbs.Uuid.String():
		return fmt.Sprintf("'%s'", gofakeit.UUID()), nil

	case dbs.Json.String(), dbs.Jsonb.String():
		fakeMap := map[string]string{
			"name":  gofakeit.FirstName(),
			"email": gofakeit.Email(),
		}
		b, _ := json.Marshal(fakeMap)
		return fmt.Sprintf("'%s'", string(b)), nil

	case dbs.Null.String():
		return "NULL", nil

	default:
		return "", fmt.Errorf("unsupported type: %s", dataType)
	}
}
