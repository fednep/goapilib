package config

import (
	"errors"
	"fmt"
	"reflect"
)

type FieldError struct {
	FieldName string
	Message   string
}

func (e FieldError) Error() string {
	return fmt.Sprintf("field %q: %s", e.FieldName, e.Message)
}

// Section defines interface used by configuration sections
type Section interface {
	IsValid() error
}

// IsValid validates Config. If any of the fields in the config struct
// implement Section interface (i.e. have IsValid() func) it is called.
//
// Error returned if any of the IsValid returns false
func IsValid(cfg any) error {

	v := reflect.ValueOf(cfg)
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return errors.New("config is nil")
		}

		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return errors.New("not a struct")
	}

	return validate(v)
}

// recursively validate config struct tree,
// ignoring any pointers to structs
func validate(st reflect.Value) error {

	t := st.Type()
	for i := 0; i < t.NumField(); i++ {
		tf := t.Field(i)
		f := st.Field(i)

		kind := tf.Type.Kind()

		// Ignore non-struct fields without tag
		if kind == reflect.Struct {
			err := validate(f)
			if err != nil {
				return fmt.Errorf("[%s] %w", tf.Name, err)
			}
		}
	}

	i, ok := st.Interface().(Section)
	if ok {
		return i.IsValid()
	}

	return nil
}
