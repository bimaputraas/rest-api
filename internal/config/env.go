package config

import (
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	JWTSecret string
	MySQLURI  string
}

func Init() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return &Config{}, err
	}
	return &Config{
		JWTSecret: os.Getenv("SECRET_JWT"),
		MySQLURI:  os.Getenv("MYSQL_URI"),
	}, nil
}
