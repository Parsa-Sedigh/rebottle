package validation

import (
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	entranslations "github.com/go-playground/validator/v10/translations/en"
	"regexp"
	"strings"
)

const IranianPhoneRegExp = "^09[0-9]{9}$"

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

	errors := make(map[string]string)

	// TODO: How we can make the field names snake case here?
	validatorErrs := err.(validator.ValidationErrors)

	// TODO: Return both fa and en translation of errors. After all that's why we used ut.Translator at the first place!
	for _, e := range validatorErrs {
		errors[strings.ToLower(e.Field())] = e.Translate(trans)
	}

	return errors
}

func ValidatePayload(validator *validator.Validate, translator ut.Translator, s interface{}) map[string]string {
	err := validator.Struct(s)

	return TranslateError(err, translator)
}
