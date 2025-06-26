package postgres

import (
	"fmt"
	"os"
	"strconv"
)

func GetDSN() (string, error) {
	host, exists := os.LookupEnv("POSTGRES_HOST")
	if !exists {
		return "", fmt.Errorf("Host Postgres is undefined")
	}
	port, exists := os.LookupEnv("POSTGRES_PORT")
	if !exists {
		return "", fmt.Errorf("Port Postgres is undefined")
	}
	user, exists := os.LookupEnv("POSTGRES_USER")
	if !exists {
		return "", fmt.Errorf("User Postgres is undefined")
	}
	password, exists := os.LookupEnv("POSTGRES_PASSWORD")
	if !exists {
		return "", fmt.Errorf("Password Postgres is undefined")
	}
	dbname, exists := os.LookupEnv("POSTGRES_DBNAME")
	if !exists {
		return "", fmt.Errorf("DBName Postgres is undefined")
	}

	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", host, user, password, dbname, port), nil
}

func GetMaxOpenConnections() (int, error) {
	maxOpenConnectionsStr, exists := os.LookupEnv("POSTGRES_MAX_OPEN_CONNECTIONS")
	if !exists {
		return 0, fmt.Errorf("POSTGRES_MAX_OPEN_CONNECTIONS environment variable not set")
	}

	maxOpenConnections, err := strconv.Atoi(maxOpenConnectionsStr)
	if err != nil {
		return 0, fmt.Errorf("POSTGRES_MAX_OPEN_CONNECTIONS must be a valid integer: %w", err)
	}

	return maxOpenConnections, nil
}

func GetMaxIdleConnections() (int, error) {
	maxIdleConnectionsStr, exists := os.LookupEnv("POSTGRES_MAX_IDLE_CONNECTIONS")
	if !exists {
		return 0, fmt.Errorf("POSTGRES_MAX_IDLE_CONNECTIONS environment variable not set")
	}

	maxIdleConnections, err := strconv.Atoi(maxIdleConnectionsStr)
	if err != nil {
		return 0, fmt.Errorf("POSTGRES_MAX_IDLE_CONNECTIONS must be a valid integer: %w", err)
	}

	return maxIdleConnections, nil
}
