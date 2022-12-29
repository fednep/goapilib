package config

import (
	"fmt"
	"os"
	"testing"
)

func TestLoadSectionFromEnv(t *testing.T) {

	type Nested struct {
		A string `env:"nested-var-a"`
	}

	type MyStruct struct {
		FieldS  string  `env:"field-string"`
		FieldPS *string `env:"field-string"`

		FieldI  int  `env:"field-int"`
		FieldPI *int `env:"field-int"`

		FieldF  float64  `env:"field-float"`
		FieldPF *float64 `env:"field-float"`

		FieldUi uint `env:"field-uint"`

		FieldB bool `env:"field-bool"`

		N     Nested  // without tag
		NP    *Nested // nested pointer (will be ignored)
		NPref Nested  `env:"n1"` // nested with tag
	}

	fStr := "some string"
	fInt := -12345
	fFloat := 1.23456

	fNestedStr := "nested string"
	fPrefixedNestedStr := "prefixed nested str"

	os.Setenv("field-string", fStr)
	os.Setenv("field-int", fmt.Sprintf("%d", fInt))
	os.Setenv("field-float", fmt.Sprintf("%f", fFloat))
	os.Setenv("field-uint", fmt.Sprintf("%d", fInt))
	os.Setenv("field-bool", "true")
	os.Setenv("nested-var-a", fNestedStr)
	os.Setenv("n1_nested-var-a", fPrefixedNestedStr)

	s := &MyStruct{}

	err := StructFromEnv(s)
	if err != nil {
		t.Fatalf("Parse struct returned error: %s", err)
	}

	// Check that pointers doesn't contain nil values
	if s.FieldPS == nil || s.FieldPI == nil || s.FieldPF == nil {
		t.Fatalf(
			"Fields: FieldPS(%v), FieldPI(%v), FieldPF(%v) should not be nil",
			s.FieldPS, s.FieldPI, s.FieldPF)
	}

	// Check for string fields
	if s.FieldS != fStr || *s.FieldPS != s.FieldS {
		t.Errorf(
			"s.FieldS(%q) and s.FieldPS(%q) contain invalid value. Expected: %q",
			s.FieldS, *s.FieldPS, fStr)
	}

	// Check for int fields
	if s.FieldI != fInt || *s.FieldPI != s.FieldI {
		t.Errorf(
			"s.FieldI(%q) and s.FieldPI(%q) contain invalid value. Expected: %q",
			s.FieldI, *s.FieldPS, fInt)
	}

	// Check for float fields
	if s.FieldF != fFloat || *s.FieldPF != s.FieldF {
		t.Errorf(
			"s.FieldF(%f) and s.FieldPF(%f) contain invalid value. Expected: %f",
			s.FieldF, *s.FieldPF, fFloat)
	}

	// Check float fields
	if s.FieldB != true {
		t.Errorf(
			"s.FieldB(%v) contain invalid value. Expected: true",
			s.FieldB)
	}

	// Check nested structure
	if s.N.A != fNestedStr {
		t.Errorf(
			"s.N.A(%v) contain invalid value. Expected: %s",
			s.N.A, fNestedStr)
	}

	// Check nested structure with struct tag prefix
	if s.NPref.A != fPrefixedNestedStr {
		t.Errorf(
			"s.NPref.A(%v) contain invalid value. Expected: %s",
			s.NPref.A, fPrefixedNestedStr)
	}
}
