package config

import (
	"github.com/joho/godotenv"
	"os"
)

// Load loads the .env file configuration
func Load() error {
	err := godotenv.Load("../.env")
	if err != nil {
		return err
	}
	return nil
}

// Get returns a particular key's value
func Get(key string) string {
	return os.Getenv(key)
}
