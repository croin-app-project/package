package http_response

import "github.com/go-playground/validator/v10"

type ErrorFieldResponse struct {
	Error       bool
	FailedField string
	Tag         string
	Value       interface{}
}

type XValidator struct {
	validator *validator.Validate
}

type GlobalErrorHandlerResp struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

var validate = validator.New()

func (v XValidator) Validate(data interface{}) []ErrorFieldResponse {
	validationErrors := []ErrorFieldResponse{}
	errs := validate.Struct(data)
	if errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			// In this case data object is actually holding the User struct
			var elem ErrorFieldResponse

			elem.FailedField = err.Field() // Export struct field name
			elem.Tag = err.Tag()           // Export struct tag
			elem.Value = err.Value()       // Export field value
			elem.Error = true

			validationErrors = append(validationErrors, elem)
		}
	}

	return validationErrors
}
