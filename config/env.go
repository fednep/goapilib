package config

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

// LoadOverrides loads data into struct from environment variables.
//
// Names of environment vars are specified using "env" structure tag.
//
// The following patterns are supported.
// Example:
//
// type Config struct {
//
//		// Standalone field with a tag
//		Var1 string `env:"SOME_VAR"`
//
//		// pointer to basic types
//		Var1 string `env:"SOME_VAR"`
//
//	 	// Nested structure without tag
//		Section ConfigSection
//
//		// Nested structure with a tag, which will be used as a prefix
//		// with underscore for any tags in the nested structure.
//		//
//		// Fore example in this case "SERVICE_" will be prefixed
//		// for any var name specified in the ConfigSection structue
//		Section2 ConfigSection `env:"SERVICE"`
//	}
func LoadOverrides(cfg any) error {
	v := reflect.ValueOf(cfg)
	if v.Kind() != reflect.Ptr {
		return errors.New("cannot load config for non-pointer")
	}

	if v.IsNil() {
		return errors.New("cannot load to nil value")
	}

	return fillStructFromEnv("", v.Elem())
}

func fillStructFromEnv(prefix string, st reflect.Value) error {
	if st.Kind() != reflect.Struct {
		return errors.New("specified interface is not a struct type")
	}

	t := st.Type()
	for i := 0; i < t.NumField(); i++ {
		tf := t.Field(i)
		f := st.Field(i)

		tag := ""
		if tf.Tag != "" {
			tag = tf.Tag.Get("env")
		}

		kind := tf.Type.Kind()

		// Ignore non-struct fields without tag
		if tag == "" && kind != reflect.Struct {
			continue
		}

		if prefix != "" {
			tag = prefix + "_" + tag
		}

		fmt.Printf("Tag: %s\n", tag)

		err := fillValue(f, tag)
		if err != nil {
			return err
		}
	}
	return nil
}

func fillValue(f reflect.Value, tag string) error {

	// TODO: Implement slice type. I.e. comma-separated list of values

	kind := f.Kind()

	switch kind {
	case reflect.Pointer:
		_, ok := os.LookupEnv(tag)
		if ok {
			if f.IsNil() {
				newVal := reflect.New(f.Type().Elem())
				err := fillValue(newVal.Elem(), tag)
				if err != nil {
					return err
				}
				f.Set(newVal)
			}
		}

	case reflect.Struct:
		err := fillStructFromEnv(tag, f)
		if err != nil {
			return err
		}

	case reflect.Float32, reflect.Float64:
		val, ok := os.LookupEnv(tag)
		if ok {
			i, err := strconv.ParseFloat(val, 64)
			if err == nil {
				f.SetFloat(i)
			}
		}

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		val, ok := os.LookupEnv(tag)
		if ok {
			i, err := strconv.ParseInt(val, 10, 64)
			if err == nil {
				f.SetInt(i)
			}
		}

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		val, ok := os.LookupEnv(tag)
		if ok {
			i, err := strconv.ParseUint(val, 10, 64)
			if err == nil {
				f.SetUint(i)
			}
		}

	case reflect.Bool:
		val, ok := os.LookupEnv(tag)
		if ok {
			vl := strings.ToLower(val)
			if vl == "false" || vl == "0" {
				f.SetBool(false)
			} else if vl == "true" || vl == "1" {
				f.SetBool(true)
			}
		}

	case reflect.String:
		val, ok := os.LookupEnv(tag)
		if ok {
			f.SetString(val)
		}
	}

	return nil
}

// LoadFromEnv fills all fields of the config structure from the environment
// variables if env:"" tags are present
func LoadFromEnv(cfg any) error {
	s := reflect.ValueOf(cfg).Elem()
	typeOfT := s.Type()
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		f.Interface()

		fmt.Printf("%d: %s %s = %v\n", i,
			typeOfT.Field(i).Name, f.Type(), f.Interface())
	}

	return nil
}

// Env loads string value from environment variable if it exists
func Env(key string, val *string) {
	if value, ok := os.LookupEnv(key); ok {
		*val = value
	}
}

// EnvInt loads int value from environment variable if it exists
func EnvInt(key string, val *int) {
	if value, ok := os.LookupEnv(key); ok {
		if v, err := strconv.Atoi(value); err == nil {
			*val = v
		}
	}
}

// EnvInt64 loads Int64 value from environment variable if it exists
func EnvInt64(key string, val *int64) {
	if value, ok := os.LookupEnv(key); ok {
		if v, err := strconv.ParseInt(value, 10, 64); err == nil {
			*val = v
		}
	}
}

// EnvBool loads bool value from environment variable if it exists
func EnvBool(key string, val *bool) {
	if value, ok := os.LookupEnv(key); ok {
		if v, err := strconv.ParseBool(value); err == nil {
			*val = v
		}
	}
}

// EnvList loads a comma-separated values from environment variable if it exists
func EnvList(key string, val *[]string) {
	if value, ok := os.LookupEnv(key); ok {
		*val = strings.Split(value, ",")
	}
}
