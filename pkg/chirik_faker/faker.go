package chirik_faker

import (
	"chickChirick/pkg/chirik_migrator/db_schema"
	"encoding/json"
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"time"
)

func FakeValue(fieldName string, dataType db_schema.DataType) string {
	switch fieldName {
	case "name":
		return gofakeit.Name()
	case "phone":
		return gofakeit.Phone()
	case "email":
		return gofakeit.Email()
	case "password":
		return gofakeit.Password(true, true, true, true, false, 20)
	case "created_at", "updated_at", "deleted_at":
		return gofakeit.TimeZoneFull()
	case "token":
		return gofakeit.UUID()
	}

	switch dataType {
	case db_schema.Bool:
		return fmt.Sprintf("%v", gofakeit.Bool())

	case db_schema.Smallint:
		return fmt.Sprintf("%d", gofakeit.IntRange(-32768, 32767))

	case db_schema.Int:
		return fmt.Sprintf("%d", gofakeit.IntRange(-2147483648, 2147483647))

	case db_schema.Bigint:
		return fmt.Sprintf("%d", gofakeit.Int64())

	case db_schema.Real:
		return fmt.Sprintf("%.4f", gofakeit.Float32Range(-1000, 1000))

	case db_schema.DoublePrecision:
		return fmt.Sprintf("%.6f", gofakeit.Float64Range(-1e6, 1e6))

	case db_schema.Varchar, db_schema.Text:
		return fmt.Sprintf("'%s'", gofakeit.Sentence(5))

	case db_schema.TimestampWithTimezone:
		return fmt.Sprintf("'%s'", gofakeit.Date().Format(time.RFC3339))

	case db_schema.TimestampWithoutTimezone:
		return fmt.Sprintf("'%s'", gofakeit.Date().Format("2006-01-02 15:04:05"))

	case db_schema.Uuid:
		return fmt.Sprintf("'%s'", gofakeit.UUID())

	case db_schema.Json, db_schema.Jsonb:
		fakeMap := map[string]string{
			"name":  gofakeit.FirstName(),
			"email": gofakeit.Email(),
		}
		b, _ := json.Marshal(fakeMap)
		return fmt.Sprintf("'%s'", string(b))

	case db_schema.Null:
		return "NULL"

	default:
		return "'unsupported_type'"
	}
}
