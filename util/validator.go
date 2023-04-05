package util

import "github.com/go-playground/validator/v10"

const (
	USD = "USD"
	EUR = "EUR"
	CAD = "CAD"
)

var ValidCurrency validator.Func = func(FieldLevel validator.FieldLevel) bool {
	if currency, ok := FieldLevel.Field().Interface().(string); ok {
		switch currency {
		case USD, EUR, CAD:
			return true
		default:
			return false
		}
	}
	return false
}
