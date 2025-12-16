package config

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func parseEnvLine(line string) (key, value string, ok bool) {
	line = strings.TrimSpace(line)

	if line == "" || strings.HasPrefix(line, "#") {
		return "", "", false
	}
	parts := strings.SplitN(line, "=", 2)
	if len(parts) != 2 {
		return "", "", false
	}

	key = strings.TrimSpace(parts[0])
	value = strings.Trim(strings.TrimSpace(parts[1]), `"'`)

	return key, value, true
}

func loadEnvVarsFromReader(r io.Reader) (map[string]string, error) {
	envMap := make(map[string]string)

	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		key, value, ok := parseEnvLine(scanner.Text())
		if ok {
			envMap[key] = value
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading content from .env: %w", err)
	}

	return envMap, nil
}

func loadEnvVarsFromFile(path string) (map[string]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("error opening .env file (%s): %w", path, err)
	}
	defer f.Close()

	return loadEnvVarsFromReader(f)
}

func LoadEnv(path string) error {
	envVars, err := loadEnvVarsFromFile(path)
	if err != nil {
		return err
	}

	for key, value := range envVars {
		if err := os.Setenv(key, value); err != nil {
			return fmt.Errorf("error setting variable %q: %w", key, err)
		}
	}

	return nil
}
