package expand

import (
	"os"
	"reflect"
)

// expandEnvStrings walks through any struct/map/slice and expands env vars in all strings.
func ExpandEnvStrings(v interface{}) {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	switch rv.Kind() {
	case reflect.String:
		rv.SetString(os.ExpandEnv(rv.String()))
	case reflect.Struct:
		for i := 0; i < rv.NumField(); i++ {
			if rv.Field(i).CanSet() || rv.Field(i).Kind() == reflect.Struct {
				ExpandEnvStrings(rv.Field(i).Addr().Interface())
			}
		}
	case reflect.Map:
		for _, key := range rv.MapKeys() {
			val := rv.MapIndex(key)
			if val.Kind() == reflect.String {
				rv.SetMapIndex(key, reflect.ValueOf(os.ExpandEnv(val.String())))
			} else {
				ExpandEnvStrings(val.Addr().Interface())
			}
		}
	case reflect.Slice:
		for i := 0; i < rv.Len(); i++ {
			ExpandEnvStrings(rv.Index(i).Addr().Interface())
		}
	}
}
