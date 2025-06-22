package domcobb

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func GetAppPort() (string, error) {
	err := godotenv.Load()
	if err != nil {
		return "", err
	}

	port, exists := os.LookupEnv("APP_PORT")
	if !exists {
		return "", fmt.Errorf("App Port is undefined!")
	}

	return port, nil
}
