package validation

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"regexp"
)

var (
	phoneRegex    = regexp.MustCompile(`^[0-9]{10}$`)
	nameRegex     = regexp.MustCompile(`^[a-zA-Z\s]+$`)
	passwordRegex = regexp.MustCompile(`^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)[A-Za-z\d]{6,20}$`)
)

func InitValidation() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("name", validateName)
		v.RegisterValidation("phone", validatePhone)
		v.RegisterValidation("password", validatePassword)
	}
}

func validateName(fl validator.FieldLevel) bool {
	return nameRegex.MatchString(fl.Field().String())
}

func validatePhone(fl validator.FieldLevel) bool {
	return phoneRegex.MatchString(fl.Field().String())
}

func validatePassword(fl validator.FieldLevel) bool {
	return passwordRegex.MatchString(fl.Field().String())
}

func FormatValidationErrors(err error) gin.H {
	var errors []string

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			switch e.Tag() {
			case "required":
				errors = append(errors, fmt.Sprintf("%s is required", e.Field()))
			case "email":
				errors = append(errors, fmt.Sprintf("%s must be a valid email", e.Field()))
			case "name":
				errors = append(errors, fmt.Sprintf("%s must contain only letters and spaces", e.Field()))
			case "phone":
				errors = append(errors, fmt.Sprintf("%s must be a valid 10-digit", e.Field()))
			case "password":
				errors = append(errors, fmt.Sprintf("%s must contain uppercase,lowercase", e.Field()))
			case "min":
				errors = append(errors, fmt.Sprintf("%s is too short", e.Field()))
			case "max":
				errors = append(errors, fmt.Sprintf("%s is too long", e.Field()))
			default:
				errors = append(errors, fmt.Sprintf("%s is invalid", e.Field()))
			}
		}
	} else {
		errors = append(errors, "Invalid request body")
	}
	return gin.H{
		"error": errors,
	}
}
