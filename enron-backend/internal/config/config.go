package config

import (
	"encoding/base64"
	"log"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	InputDir          string
	ZincURL           string
	IndexName         string
	BasicAuthToken    string
	FoldersPerRoutine int
}

func LoadConfig() (*Config, error) {
	if err := LoadEnv(); err != nil {
		log.Println("No se encontró archivo .env, se utilizarán las variables de entorno")
	}

	foldersPerRoutine, err := strconv.Atoi(os.Getenv("FOLDERS_PER_ROUTINE"))
	if err != nil {
		foldersPerRoutine = 15
	}

	cfg := &Config{
		InputDir:          os.Getenv("INPUT_DIR"),
		ZincURL:           os.Getenv("ZINC_URL"),
		IndexName:         os.Getenv("INDEX_NAME"),
		BasicAuthToken:    "Basic " + base64.StdEncoding.EncodeToString([]byte(os.Getenv("ZINC_FIRST_ADMIN_USER")+":"+os.Getenv("ZINC_FIRST_ADMIN_PASS"))),
		FoldersPerRoutine: foldersPerRoutine,
	}

	return cfg, nil
}

func LoadEnv() error {
	data, err := os.ReadFile(".env")
	if err != nil {
		return err
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key, value := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
		os.Setenv(key, value)
	}

	return nil
}
