package middleware

import "net/http"

// ValidatorInterface TODO: текущая реализация валидатора не подразумевает возврат ошибки на уровень контроллера. Доработать
type ValidatorInterface interface {
	Validate(next http.HandlerFunc) http.HandlerFunc
}
