package helpers

import (
	"errors"
	enLocalePackage "github.com/go-playground/locales/en"
	universalTranslatorPackage "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslationsPackage "github.com/go-playground/validator/v10/translations/en"
)

var (
	universal *universalTranslatorPackage.UniversalTranslator
	validate  *validator.Validate
)

func Validate(validateStruct interface{}) (validateErrs []string, err error) {
	en := enLocalePackage.New()
	universal = universalTranslatorPackage.New(en, en)
	translator, localeFound := universal.GetTranslator("en")
	if !localeFound {
		return nil, errors.New("translator locale not found")
	}

	validate = validator.New()
	if err = enTranslationsPackage.RegisterDefaultTranslations(validate, translator); err != nil {
		return
	}

	structErrs := validate.Struct(validateStruct)
	if structErrs != nil {
		errs := structErrs.(validator.ValidationErrors).Translate(translator)

		for _, err := range errs {
			validateErrs = append(validateErrs, err)
		}
	}

	return
}
