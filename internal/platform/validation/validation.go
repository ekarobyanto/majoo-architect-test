package validation

import (
	"github.com/go-playground/validator/v10"
	"github.com/user/simple-blog/internal/platform/errors"
)

var validate = validator.New()

type ErrorField struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
	Value string `json:"value,omitempty"`
}

func Validate(s interface{}) error {
	err := validate.Struct(s)
	if err != nil {
		var details []ErrorField
		for _, err := range err.(validator.ValidationErrors) {
			details = append(details, ErrorField{
				Field: err.Field(),
				Tag:   err.Tag(),
				Value: err.Param(),
			})
		}
		return errors.ValidationError("Validation failed", details)
	}
	return nil
}
