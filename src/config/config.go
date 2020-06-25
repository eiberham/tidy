package config

import (
	"fmt"
	"github.com/joho/godotenv"
)

func Load() {
	err := godotenv.Load("../.env")
	if err != nil {
		fmt.Println("Error loading the configuration")
	}
}
