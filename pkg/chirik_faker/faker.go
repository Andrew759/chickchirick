package chirik_faker

import (
	dbs "chickChirick/pkg/chirik_migrator/db_schema"
	"encoding/json"
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"time"
)

func FakeValue(fieldName string, dataType string) string {
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
	case string(dbs.Bool):
		return fmt.Sprintf("%v", gofakeit.Bool())

	case string(dbs.Smallint):
		return fmt.Sprintf("%d", gofakeit.IntRange(-32768, 32767))

	case string(dbs.Int):
		return fmt.Sprintf("%d", gofakeit.IntRange(-2147483648, 2147483647))

	case string(dbs.Bigint):
		return fmt.Sprintf("%d", gofakeit.Int64())

	case string(dbs.BigSerial):
		return fmt.Sprintf("%d", gofakeit.Int64())

	case string(dbs.Real):
		return fmt.Sprintf("%.4f", gofakeit.Float32Range(-1000, 1000))

	case string(dbs.DoublePrecision):
		return fmt.Sprintf("%.6f", gofakeit.Float64Range(-1e6, 1e6))

	case string(dbs.Varchar), string(dbs.Text):
		return fmt.Sprintf("'%s'", gofakeit.Sentence(5))

	case string(dbs.TimestampWithTimezone):
		return fmt.Sprintf("'%s'", gofakeit.Date().Format(time.RFC3339))

	case string(dbs.TimestampWithoutTimezone):
		return fmt.Sprintf("'%s'", gofakeit.Date().Format("2006-01-02 15:04:05"))

	case string(dbs.Uuid):
		return fmt.Sprintf("'%s'", gofakeit.UUID())

	case string(dbs.Json), string(dbs.Jsonb):
		fakeMap := map[string]string{
			"name":  gofakeit.FirstName(),
			"email": gofakeit.Email(),
		}
		b, _ := json.Marshal(fakeMap)
		return fmt.Sprintf("'%s'", string(b))

	case string(dbs.Null):
		return "NULL"

	default:
		return "'unsupported_type'"
	}
}
