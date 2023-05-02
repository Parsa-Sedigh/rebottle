package validation

import (
	"fmt"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"regexp"
)

const IranianPhoneRegExp = "/^09[0-9]{9}$/"

func Register() (*validator.Validate, error) {
	validate := validator.New()
	err := validate.RegisterValidation("phone", ValidatePhone)
	if err != nil {
		return nil, err
	}

	return validate, nil
}

func ValidatePhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()

	return regexp.MustCompile(IranianPhoneRegExp).MatchString(phone)
}

func TranslateError(err error, trans ut.Translator) (errs []error) {
	if err == nil {
		return nil
	}

	validatorErrs := err.(validator.ValidationErrors)
	for _, e := range validatorErrs {
		translatedErr := fmt.Errorf(e.Translate(trans))
		errs = append(errs, translatedErr)
	}

	return errs
}
