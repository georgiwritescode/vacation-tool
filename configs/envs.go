package configs

import (
	"os"
)

type Config struct {
	PublicHost string
	Port       string
	DBUser     string
	DBPassword string
	DBAddress  string
	DBName     string
}

var Envs = initConfigs()

func initConfigs() Config {

	return Config{
		PublicHost: getEnv("HOST", "localhost"),
		Port:       getEnv("PORT", ":8080"),
		DBUser:     getEnv("DB_USER", "portal"),
		DBPassword: getEnv("DB_PASSWORD", "password123"),
		DBAddress:  getEnv("DB_ADDRESS", "127.0.0.1:3307"),
		DBName:     getEnv("DB_NAME", "vacation_tool"),
	}
}

func getEnv(key, fallback string) string {
	if value, err := os.LookupEnv(key); err {
		return value
	}

	return fallback
}
