package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// LoadEnv lê um arquivo .env e coloca cada variável no ambiente do processo
// usando os.Setenv(KEY, VALUE).
func LoadEnv(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("erro abrindo arquivo .env (%s): %w", path, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Ignora linha vazia ou comentário
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Formato esperado: CHAVE=VALOR
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			// Linha mal formatada, ignora ou loga se quiser
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Remove aspas se estiverem presentes
		value = strings.Trim(value, `"'`)

		os.Setenv(key, value)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("erro lendo arquivo .env: %w", err)
	}

	return nil
}
