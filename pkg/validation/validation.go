package validation

import (
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	entranslations "github.com/go-playground/validator/v10/translations/en"
	"regexp"
)

const IranianPhoneRegExp = "/^09[0-9]{9}$/"

func Register() (*validator.Validate, ut.Translator, error) {
	validate := validator.New()
	english := en.New()
	uni := ut.New(english, english)
	trans, _ := uni.GetTranslator("en")
	err := entranslations.RegisterDefaultTranslations(validate, trans)
	if err != nil {
		return nil, nil, err
	}

	err = validate.RegisterValidation("phone", ValidatePhone)
	if err != nil {
		return nil, nil, err
	}

	return validate, trans, nil
}

func ValidatePhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()

	return regexp.MustCompile(IranianPhoneRegExp).MatchString(phone)
}

/* TranslateError creates a map of errors with the key as the name of the field and value as the actual error  */
func TranslateError(err error, trans ut.Translator) map[string]string {
	if err == nil {
		return nil
	}

	errs := make(map[string]string)

	validatorErrs := err.(validator.ValidationErrors)

	for _, e := range validatorErrs {
		errs[e.Field()] = e.Translate(trans)
	}

	return errs
}
