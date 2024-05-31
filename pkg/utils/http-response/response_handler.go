package http_response

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

type SuccessReponse struct {
	Code   int         `json:"code"`   // Response code
	Status string      `json:"status"` // Response status
	Result interface{} `json:"result"` // Response result
}

type ErrorResponse struct {
	Code    int         `json:"code"`    // Response error code
	Message string      `json:"message"` // Response error message
	Error   interface{} `json:"error"`   // Response error detail
}

const (
	SUCCESS = 200

	INVALID_INPUT_PARAMETER = -40001
	DATA_NOT_FOUND          = -40002
	DATA_ALREADY_EXISTS     = -40003

	INVALID_AUTHORIZATION_KEY  = -40101
	AUTHORIZATION_KEY_EXPIRED  = -40102
	AUTHORIZATION_KEY_INACTIVE = -40103
	USER_NOT_FOUND             = -40104
	INVALID_PASSWORD           = -40105
	USER_EXPIRED               = -40106
	PASSWORD_RESETED           = -40107

	INTERNAL_SYSTEM_ERROR = -50001
	DATABASE_ERROR        = -50002
	NOT_FOUND             = -50003
)

var ResponseMessages = map[int]string{
	SUCCESS: "Success",

	INVALID_INPUT_PARAMETER: "Invalid input parameter",
	DATA_NOT_FOUND:          "Data Not Found",
	DATA_ALREADY_EXISTS:     "Data already exists",

	INVALID_AUTHORIZATION_KEY:  "Invalid authorization key",
	AUTHORIZATION_KEY_EXPIRED:  "Authorization key expired",
	AUTHORIZATION_KEY_INACTIVE: "Authorization key inactive",
	USER_NOT_FOUND:             "User Not Found",
	INVALID_PASSWORD:           "Invalid Password",
	USER_EXPIRED:               "User Expired",
	PASSWORD_RESETED:           "Password was reset by admin please change password via WMS on Windows (แจ้งให้User ไปท าการ Change Password ที่ WMS on Windows)",

	INTERNAL_SYSTEM_ERROR: "Internal system error",
	DATABASE_ERROR:        "Database error",
	NOT_FOUND:             "Not Found",
}

func HandleException(errorCode int, err error) (int, ErrorResponse) {
	message, ok := ResponseMessages[errorCode]
	if !ok {
		message = fmt.Sprintf("%s (%d)", ResponseMessages[INTERNAL_SYSTEM_ERROR], errorCode)
	}

	var errors interface{}
	if err != nil {
		errors = parseError(err)
	} else {
		errors = nil
	}

	if errorCode == INVALID_INPUT_PARAMETER {
		return http.StatusBadRequest, ErrorResponse{
			Code:    errorCode,
			Message: message,
			Error:   errors,
		}
	} else if errorCode == DATA_NOT_FOUND {
		return http.StatusBadRequest, ErrorResponse{
			Code:    errorCode,
			Message: message,
			Error:   errors,
		}
	} else if errorCode == DATA_ALREADY_EXISTS {
		return http.StatusBadRequest, ErrorResponse{
			Code:    errorCode,
			Message: message,
			Error:   errors,
		}
	} else if errorCode == INVALID_AUTHORIZATION_KEY {
		return http.StatusUnauthorized, ErrorResponse{
			Code:    errorCode,
			Message: message,
			Error:   errors,
		}
	} else if errorCode == AUTHORIZATION_KEY_EXPIRED {
		return http.StatusUnauthorized, ErrorResponse{
			Code:    errorCode,
			Message: message,
			Error:   errors,
		}
	} else if errorCode == AUTHORIZATION_KEY_INACTIVE {
		return http.StatusUnauthorized, ErrorResponse{
			Code:    errorCode,
			Message: message,
			Error:   errors,
		}
	} else if errorCode == USER_NOT_FOUND {
		return http.StatusUnauthorized, ErrorResponse{
			Code:    errorCode,
			Message: message,
			Error:   errors,
		}
	} else if errorCode == INVALID_PASSWORD {
		return http.StatusUnauthorized, ErrorResponse{
			Code:    errorCode,
			Message: message,
			Error:   errors,
		}
	} else if errorCode == USER_EXPIRED {
		return http.StatusUnauthorized, ErrorResponse{
			Code:    errorCode,
			Message: message,
			Error:   errors,
		}
	} else if errorCode == PASSWORD_RESETED {
		return http.StatusUnauthorized, ErrorResponse{
			Code:    errorCode,
			Message: message,
			Error:   errors,
		}
	} else if errorCode == DATABASE_ERROR {
		return http.StatusInternalServerError, ErrorResponse{
			Code:    errorCode,
			Message: message,
			Error:   errors,
		}
	} else if errorCode == NOT_FOUND {
		return http.StatusInternalServerError, ErrorResponse{
			Code:    errorCode,
			Message: message,
			Error:   errors,
		}
	} else {
		return http.StatusInternalServerError, ErrorResponse{
			Code:    errorCode,
			Message: message,
			Error:   errors,
		}
	}
}

func parseError(errs ...error) []string {
	var out []string
	for _, err := range errs {
		switch typedError := any(err).(type) {
		case validator.ValidationErrors:
			// if the type is validator.ValidationErrors then it's actually an array of validator.FieldError so we'll
			// loop through each of those and convert them one by one
			for _, e := range typedError {
				out = append(out, parseFieldError(e))
			}
		case *json.UnmarshalTypeError:
			// similarly, if the error is an unmarshalling error we'll parse it into another, more readable string format
			out = append(out, parseMarshallingError(*typedError))
		default:
			out = append(out, err.Error())
		}
	}
	return out
}
func parseFieldError(e validator.FieldError) string {
	// workaround to the fact that the `gt|gtfield=Start` gets passed as an entire tag for some reason
	// https://github.com/go-playground/validator/issues/926
	fieldPrefix := fmt.Sprintf("The field %s", e.Field())
	tag := strings.Split(e.Tag(), "|")[0]
	switch tag {
	case "required_without":
		return fmt.Sprintf("%s is required if %s is not supplied", fieldPrefix, e.Param())
	case "lt", "ltfield":
		param := e.Param()
		if param == "" {
			param = time.Now().Format(time.RFC3339)
		}
		return fmt.Sprintf("%s must be less than %s", fieldPrefix, param)
	case "gt", "gtfield":
		param := e.Param()
		if param == "" {
			param = time.Now().Format(time.RFC3339)
		}
		return fmt.Sprintf("%s must be greater than %s", fieldPrefix, param)
	default:
		// if it's a tag for which we don't have a good format string yet we'll try using the default english translator
		english := en.New()
		translator := ut.New(english, english)
		if translatorInstance, found := translator.GetTranslator("en"); found {
			return e.Translate(translatorInstance)
		} else {
			return fmt.Errorf("%v", e).Error()
		}
	}
}
func parseMarshallingError(e json.UnmarshalTypeError) string {
	return fmt.Sprintf("The field %s must be a %s", e.Field, e.Type.String())
}
