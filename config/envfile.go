package config

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

// LoadEnv loads .env file from the current working directory if it exists
// envFile is just an ".env" or ".env.dev" without full file path
//
// To load using full file path, use LoadEnvFile
func LoadEnv(envFile string) error {

	if envFile == "" {
		envFile = ".env"
	}

	fn, err := EnvFile(envFile)
	if err != nil {
		return fmt.Errorf("Cannot get .env file name: %s", err)
	}

	if fn == "" {
		return nil
	}

	return LoadEnvFile(fn)
}

// EnvFile returns full path to .env file if it exists in the current working
// directory. If file doesn't exists, empty string returned
//
// suffix is used to append to the .env file,
// for example for "dev" suffix, ".env.dev" will be looked up
func EnvFile(fn string) (string, error) {
	pwd, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return "", err
	}

	fn = path.Join(pwd, fn)
	exists, err := fileExists(fn)
	if err != nil {
		return "", err
	}

	if !exists {
		return "", nil
	}

	return fn, nil
}

// LoadEnvFile reads all strings from fn file
// and loads it into the environment
// The format of the string is key=value
func LoadEnvFile(fn string) error {
	log.Printf("Loading env vars from %q file", fn)
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
