package config

import (
	"testing"
)

func TestParseEnvLine(t *testing.T) {

	testCases := []string{
		"somkey = this is a string ",
		"somkey=this is a string",
		"somekey=\"this is a string\"",
		"somekey = \"this is a string\"",
	}

	for _, c := range testCases {
		key, value, err := parseEnvLine(c)
		if err != nil {
			t.Fatalf("Error parsing use case %q: %s", c, err)
		}
		if key != "somekey" && value != "this is a string" {
			t.Fatalf("Error parsing for the case %q: key (%q) or value (%q) incorrect", c, key, value)
		}
	}
}
