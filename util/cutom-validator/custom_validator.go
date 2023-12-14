package customvalidator

import (
	"github.com/aryyawijaya/simple-bank/util"
	"github.com/go-playground/validator/v10"
)

var ValidCurrency validator.Func = func(fl validator.FieldLevel) bool {
	if currency, ok := fl.Field().Interface().(string); ok {
		// validate currency
		return util.IsSupportedCurrency(currency)
	}

	return false
}
