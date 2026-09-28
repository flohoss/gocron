package validate

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

var (
	validate = validator.New()
	trans    ut.Translator
)

func init() {
	translator := ut.New(en.New())
	trans, _ = translator.GetTranslator("en")

	validate.RegisterTagNameFunc(configFieldName)
	if err := en_translations.RegisterDefaultTranslations(validate, trans); err != nil {
		panic(err)
	}

	registerRules()
}

func Struct(target any) error {
	return pretty(validate.Struct(target))
}

func Var(value any, rules string) error {
	return pretty(validate.Var(value, rules))
}

func RegisterValidation(tag string, fn validator.Func) error {
	return validate.RegisterValidation(tag, fn)
}

// RegisterRule registers a rule and the message the pretty printer renders for
// it. Rules registered without a translation fall back to the unhelpful
// "field failed the \"tag\" rule".
func RegisterRule(tag, message string, fn validator.Func) error {
	if err := validate.RegisterValidation(tag, fn); err != nil {
		return err
	}

	return validate.RegisterTranslation(tag, trans, func(ut ut.Translator) error {
		return ut.Add(tag, message, true)
	}, func(ut ut.Translator, fieldError validator.FieldError) string {
		translated, err := ut.T(tag, fieldError.Field(), fieldError.Param())
		if err != nil {
			return fieldError.Error()
		}

		return translated
	})
}

func pretty(err error) error {
	var fieldErrors validator.ValidationErrors
	if !errors.As(err, &fieldErrors) {
		return err
	}

	lines := make([]string, 0, len(fieldErrors))
	for _, fieldError := range fieldErrors {
		lines = append(lines, "- "+formatFieldError(fieldError))
	}

	return errors.New(strings.Join(lines, "\n"))
}

func formatFieldError(fieldError validator.FieldError) string {
	path := configPath(fieldError.Namespace())
	translated := fieldError.Translate(trans)

	if translated == fieldError.Error() {
		return fmt.Sprintf("%s failed the %q rule", path, fieldError.Tag())
	}

	return path + strings.TrimPrefix(translated, fieldError.Field())
}

func configPath(namespace string) string {
	parts := strings.Split(namespace, ".")
	if len(parts) > 1 {
		parts = parts[1:]
	}

	return strings.Join(parts, ".")
}

func configFieldName(field reflect.StructField) string {
	name, _, _ := strings.Cut(field.Tag.Get("mapstructure"), ",")
	if name == "" || name == "-" {
		return camelCase(field.Name)
	}

	return name
}

func camelCase(name string) string {
	if name == "" {
		return name
	}

	return strings.ToLower(name[:1]) + name[1:]
}
