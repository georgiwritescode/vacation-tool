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
		DBUser:     getEnv("DB_User", "portal"),
		DBPassword: getEnv("DB_Password", "password123"),
		DBAddress:  getEnv("DB_Address", "127.0.0.1:3307"),
		DBName:     getEnv("DB_Name", "vacation_tool"),
	}
}

func getEnv(key, fallback string) string {
	if value, err := os.LookupEnv(key); err {
		return value
	}

	return fallback
}
