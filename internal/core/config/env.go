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
	value = strings.TrimSpace(parts[1])

	// Remove aspas simples/duplas se estiverem presentes nas extremidades
	value = strings.Trim(value, `"'`)

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
		return nil, fmt.Errorf("Error reading content from .env: %w", err)
	}

	return envMap, nil
}

func loadEnvVarsFromFile(path string) (map[string]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("Error opening .env file (%s): %w", path, err)
	}
	defer f.Close()

	return loadEnvVarsFromReader(f)
}

// LoadEnv carrega as variáveis do arquivo e aplica no ambiente do processo.
func LoadEnv(path string) error {
	envVars, err := loadEnvVarsFromFile(path)
	if err != nil {
		return err
	}

	for key, value := range envVars {
		if err := os.Setenv(key, value); err != nil {
			return fmt.Errorf("Error setting variable %q: %w", key, err)
		}
	}

	return nil
}
