package router

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/go-playground/validator"
)

type CustomValidator struct {
	Validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.Validator.Struct(i)
}

func newValidator() *validator.Validate {
	validate := validator.New()
	validate.RegisterValidation(`rfg`, func(fl validator.FieldLevel) bool {
		param := strings.Split(fl.Param(), `:`)
		paramField := param[0]
		paramValue := param[1]
		if paramField == `` {
			return true
		}
		var paramFieldValue reflect.Value
		if fl.Parent().Kind() == reflect.Ptr {
			paramFieldValue = fl.Parent().Elem().FieldByName(paramField)
		} else {
			paramFieldValue = fl.Parent().FieldByName(paramField)
		}
		if !isGt(paramFieldValue, paramValue) {
			return true
		}
		return hasValue(fl)
	})
	return validate
}
func hasValue(fl validator.FieldLevel) bool {
	return requireCheckFieldKind(fl, "")
}

func requireCheckFieldKind(fl validator.FieldLevel, param string) bool {
	field := fl.Field()
	if len(param) > 0 {
		if fl.Parent().Kind() == reflect.Ptr {
			field = fl.Parent().Elem().FieldByName(param)
		} else {
			field = fl.Parent().FieldByName(param)
		}
	}
	switch field.Kind() {
	case reflect.Slice, reflect.Map, reflect.Ptr, reflect.Interface, reflect.Chan, reflect.Func:
		return !field.IsNil()
	default:
		_, _, nullable := fl.ExtractType(field)
		if nullable && field.Interface() != nil {
			return true
		}
		return field.IsValid() && field.Interface() != reflect.Zero(field.Type()).Interface()
	}
}

func isGt(field reflect.Value, value string) bool {
	switch field.Kind() {

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p := asInt(value)

		return field.Int() > p

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		p := asUint(value)

		return field.Uint() > p

	case reflect.Float32, reflect.Float64:
		p := asFloat(value)

		return field.Float() > p
	}

	panic(fmt.Sprintf("Bad field type %T", field.Interface()))
}

func asInt(param string) int64 {

	i, err := strconv.ParseInt(param, 0, 64)
	panicIf(err)

	return i
}

func asUint(param string) uint64 {

	i, err := strconv.ParseUint(param, 0, 64)
	panicIf(err)

	return i
}

func asFloat(param string) float64 {

	i, err := strconv.ParseFloat(param, 64)
	panicIf(err)

	return i
}

func panicIf(err error) {
	if err != nil {
		panic(err.Error())
	}
}
