package configs

import "github.com/joho/godotenv"

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
	godotenv.Load()

	//todo: implement some sort of reading from .env file
	return Config{
		PublicHost: "http://localhost",
		Port:       ":8080",
		DBUser:     "portal",
		DBPassword: "password123",
		DBAddress:  "127.0.0.1:3306",
		DBName:     "vacation_tool",
	}
}
