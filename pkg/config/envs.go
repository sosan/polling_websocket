package config

import (
	"log"
	"os"
	"path/filepath"
	"polling_websocket/pkg/vaults"

	"github.com/joho/godotenv"
)

func LoadEnvs(baseDir string) {
	if err := loadCurrentEnv(); err != nil {
		log.Printf("WARNING | Cannot read current .env: %v", err)
		if err := loadBaseEnv(baseDir, ".env"); err != nil {
			log.Printf("WARNING | Cannot read base .env: %v", err)
			if err := loadBaseEnv(baseDir, ".env.local"); err != nil {
				log.Printf("WARNING | Cannot read local .env.local: %v", err)
			}
		}
	}

	// Load environment variables from Vault
	LoadEnvsFromVault()
}

func loadCurrentEnv() error {
	return godotenv.Load()
}

func loadBaseEnv(baseDir string, fileName string) error {
	envPath := filepath.Join(baseDir, fileName)
	return loadEnvFile(envPath)
}

func loadEnvFile(envFilePath string) error {
	if _, err := os.Stat(envFilePath); err != nil {
		return err
	}
	return godotenv.Load(envFilePath)
}

func LoadEnvsFromVault() {
	vaults.GetEnvsFromVault()
}

func GetEnv(key, fallback string) string {
	if value, exist := os.LookupEnv(key); exist {
		return value
	}
	return fallback
}
