package middleware

import (
	mainService "chickChirick/cmd/service"
	"net/http"
)

type ValidatorFactory interface {
	NewValidator(dbDecorator mainService.DBDecorator, opts ...ValidatorOption) Validator
}

type Validator interface {
	Validate(next http.HandlerFunc) http.HandlerFunc
}

type ValidatorOptions struct { //Конфигурация структуры
	validateDB bool
}

type ValidatorOption func(options *ValidatorOptions)

func ActivateDBValidation() ValidatorOption { //Функция конфигурации,
	return func(vOptions *ValidatorOptions) {
		vOptions.validateDB = true
	}
}

func (vo ValidatorOptions) IsDBValidationActivated() bool {
	return vo.validateDB
}
