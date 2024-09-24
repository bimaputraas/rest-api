package pkgvalidate

import (
	"errors"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

// use a single instance of Validate, it caches struct info
var validate *validator.Validate = validator.New(validator.WithRequiredStructEnabled())

func Struct(obj interface{}) error {
	err := validate.Struct(obj)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return errors.New("failed assert *validator.InvalidValidationError")
		}

		for _, err := range err.(validator.ValidationErrors) {
			switch err.Tag() {
			case "email":
				return errors.New("invalid email")
			case "required":
				return errors.New(strings.ToLower(err.Field()) + " is required")
			case "e164":
				return errors.New("phone is not in E.164 format")
			default:
				return errors.New("invalid " + err.Field())
			}
		}
	}

	return nil
}

func IsNumeric(s string) bool {
	re := regexp.MustCompile(`^[0-9]+$`)
	return re.MatchString(s)
}
