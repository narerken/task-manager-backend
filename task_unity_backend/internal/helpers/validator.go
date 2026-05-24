package helpers

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

func NewValidator() *validator.Validate {
	v := validator.New()
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	return v
}

func Validate(req interface{}, v *validator.Validate) map[string][]string {
	if v == nil {
		return map[string][]string{"error": {"validator is not initialized"}}
	}

	errors := make(map[string][]string)

	switch reflect.TypeOf(req).Kind() {
	case reflect.Slice, reflect.Array:
		s := reflect.ValueOf(req)
		for i := 0; i < s.Len(); i++ {
			item := s.Index(i).Interface()
			errs := v.Struct(item)
			if errs != nil {
				for _, err := range errs.(validator.ValidationErrors) {
					fieldKey := fmt.Sprintf("items.%d.%s", i, err.Field())
					errors[fieldKey] = append(errors[fieldKey], getErrorMessage(err))
				}
			}
		}
	default:
		errs := v.Struct(req)
		if errs != nil {
			for _, err := range errs.(validator.ValidationErrors) {
				errors[err.Field()] = append(errors[err.Field()], getErrorMessage(err))
			}
		}
	}

	if len(errors) > 0 {
		return errors
	}
	return nil
}

func getErrorMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return fmt.Sprintf("Field '%s' is required", err.Field())
	case "min":
		return fmt.Sprintf("Field '%s' must be at least %s characters", err.Field(), err.Param())
	case "max":
		return fmt.Sprintf("Field '%s' must be at most %s characters", err.Field(), err.Param())
	case "email":
		return fmt.Sprintf("Field '%s' must be a valid email", err.Field())
	case "oneof":
		return fmt.Sprintf("Field '%s' must be one of [%s]", err.Field(), err.Param())
	default:
		return fmt.Sprintf("Field '%s' is invalid", err.Field())
	}
}
