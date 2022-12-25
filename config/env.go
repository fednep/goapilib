package config

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// LoadDotEnv loads .env file from the current working directory
func LoadDotEnv() error {

	// Get current directory
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return fmt.Errorf("Cannot get absolute path of executable: %s", err)
	}

	return LoadEnvFile(path.Join(dir, ".env"))
}

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
