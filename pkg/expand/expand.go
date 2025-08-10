package expand

import (
	"os"
	"reflect"
)

func ExpandEnvStrings(v any) {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	switch rv.Kind() {
	case reflect.String:
		rv.SetString(os.ExpandEnv(rv.String()))

	case reflect.Struct:
		for i := 0; i < rv.NumField(); i++ {
			field := rv.Field(i)
			if field.CanSet() && field.Kind() == reflect.String {
				field.SetString(os.ExpandEnv(field.String()))
			} else if field.Kind() == reflect.Ptr || field.Kind() == reflect.Struct || field.Kind() == reflect.Map || field.Kind() == reflect.Slice {
				if field.CanAddr() {
					ExpandEnvStrings(field.Addr().Interface())
				}
			}
		}

	case reflect.Map:
		for _, key := range rv.MapKeys() {
			val := rv.MapIndex(key)

			switch val.Kind() {
			case reflect.String:
				newVal := os.ExpandEnv(val.String())
				rv.SetMapIndex(key, reflect.ValueOf(newVal))
			case reflect.Ptr, reflect.Struct, reflect.Map, reflect.Slice:
				if val.CanAddr() {
					ExpandEnvStrings(val.Addr().Interface())
				}
			}
		}

	case reflect.Slice:
		for i := 0; i < rv.Len(); i++ {
			elem := rv.Index(i)
			if elem.Kind() == reflect.String && elem.CanSet() {
				elem.SetString(os.ExpandEnv(elem.String()))
			} else if elem.Kind() == reflect.Ptr || elem.Kind() == reflect.Struct || elem.Kind() == reflect.Map || elem.Kind() == reflect.Slice {
				if elem.CanAddr() {
					ExpandEnvStrings(elem.Addr().Interface())
				}
			}
		}
	}
}
