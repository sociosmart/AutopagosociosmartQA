package utils

import (
	"reflect"
	"smartgas-payment/internal/dto"

	"github.com/go-playground/validator/v10"
)

// MapValidatorError takes a generic parameter in order to get the json tag
func MapValidatorError[T any](errs error) []*dto.BadRequestMessage {
	response := make([]*dto.BadRequestMessage, 0)

	for _, err := range errs.(validator.ValidationErrors) {
		field, _ := reflect.TypeOf(new(T)).Elem().FieldByName(err.Field())
		jsonTag, _ := field.Tag.Lookup("json")
		response = append(response, &dto.BadRequestMessage{
			Field:   jsonTag,
			Message: err.Tag(),
		})
	}

	return response
}
