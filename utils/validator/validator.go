package validator

import (
	"regexp"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func InitValidator() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// 1. 注册密码验证
		v.RegisterValidation("password_rule", passwordValidator)
	}
}

func passwordValidator(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	if len(password) < 6 || len(password) > 20 {
		return false
	}
	hasDigit, _ := regexp.MatchString(`[0-9]`, password)
	hasUpper, _ := regexp.MatchString(`[A-Z]`, password)
	hasLower, _ := regexp.MatchString(`[a-z]`, password)
	return hasDigit && hasUpper && hasLower
}
