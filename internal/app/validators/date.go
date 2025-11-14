package validators

import (
	"reflect"
	"regexp"

	"github.com/go-playground/validator/v10"
)

func DateValidator() validator.Func {
	re := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)

	return func(fl validator.FieldLevel) bool {
		field := fl.Field()

		if field.Kind() == reflect.String {
			return re.MatchString(field.String())
		} else {
			return false
		}
	}
}
