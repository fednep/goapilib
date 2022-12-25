package config

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// LoadEnvFile reads all strings from fn file
// and loads it into the environment
// The format of the string is key=value
func LoadEnvFile(fn string) error {
	f, err := os.Open(fn)
	if err != nil {
		return fmt.Errorf("cannot open file: %q: %w", fn, err)
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		key, value, err := parseEnvLine(scanner.Text())
		if err != nil {
			continue
		}

		err = os.Setenv(key, value)
		if err != nil {
			return fmt.Errorf("cannot set env variable: %w", err)
		}
	}

	return nil
}

// parseEnvLine splits line into the key, value pairs. Line can be in the
// following forms:
//
//	key=value
//	key="value"
//
// When quotes are used, they are not returned as part of the value.
func parseEnvLine(line string) (string, string, error) {
	re := regexp.MustCompile(`^(\w+)\s*=\s*(?:"(.*)"|(.*?))\s*$`)
	if strings.HasPrefix(line, "#") || !re.MatchString(line) {
		return "", "", fmt.Errorf("cannot parse line. No key/value pair found")
	}

	sm := re.FindStringSubmatch(line)

	// If value specified in quotes (second group)
	if sm[2] != "" {
		return sm[1], sm[2], nil
	}

	return sm[1], sm[3], nil
}
