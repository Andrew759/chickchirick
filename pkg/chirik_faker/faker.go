package chirik_faker

import (
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"time"
)

// FakeValue TODO: сейчас нет обработки json
func FakeValue(fieldName any) (string, error) {
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

	switch fieldName.(type) {
	case bool:
		return fmt.Sprintf("%v", gofakeit.Bool()), nil
	case int8:
		return fmt.Sprintf("%d", gofakeit.Int8()), nil
	case int16:
		return fmt.Sprintf("%d", gofakeit.Int16()), nil
	case int32:
		return fmt.Sprintf("%d", gofakeit.Int32()), nil
	case int64:
		return fmt.Sprintf("%d", gofakeit.Int64()), nil
	case int:
		return fmt.Sprintf("%d", gofakeit.Int64()), nil
	case uint8:
		return fmt.Sprintf("%d", gofakeit.Uint8()), nil
	case uint16:
		return fmt.Sprintf("%d", gofakeit.Uint16()), nil
	case uint32:
		return fmt.Sprintf("%d", gofakeit.Uint32()), nil
	case uint64:
		return fmt.Sprintf("%d", gofakeit.Uint64()), nil
	case uint:
		return fmt.Sprintf("%d", gofakeit.Uint64()), nil
	case float32:
		return fmt.Sprintf("%f", gofakeit.Float32()), nil
	case float64:
		return fmt.Sprintf("%f", gofakeit.Float64()), nil
	case string:
		return fmt.Sprintf("'%s'", gofakeit.Sentence(5)), nil
	case time.Time:
		return fmt.Sprintf("'%s'", gofakeit.Date().Format(time.RFC3339)), nil
	default:
		return "", fmt.Errorf("unsupported type")
	}
}
