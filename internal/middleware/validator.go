package middleware

import "net/http"

type ValidatorInterface interface {
	Validate(next http.HandlerFunc) http.HandlerFunc
}
