package vaults

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func GetEnvsFromVault() {
	envValues := GetAllEnvsFromRedis()
	envsMaped := serializeEnvs(envValues)
	SetEnvs(envsMaped)
}

func serializeEnvs(envStr string) map[string]string {
	envMap, err := godotenv.Unmarshal(envStr)
	if err != nil {
		log.Panic("ERROR | Cannot serialize string from Env")
	}
	return envMap
}

func SetEnvs(envsMapped map[string]string) {
	for key, value := range envsMapped {
		os.Setenv(key, value)
	}
}
