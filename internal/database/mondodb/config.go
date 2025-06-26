package mondodb

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

func GetURI() (string, error) {
	host, exists := os.LookupEnv("MONGO_HOST")
	if !exists {
		return "", fmt.Errorf("Mongo host is undefined")
	}
	port, exists := os.LookupEnv("MONGO_PORT")
	if !exists {
		return "", fmt.Errorf("Mongo Port is endefined")
	}
	username, exists := os.LookupEnv("MONGO_ROOT_USERNAME")
	if !exists {
		return "", fmt.Errorf("Mongo username is undefined")
	}
	password, exists := os.LookupEnv("MONGO_ROOT_PASSWORD")
	if !exists {
		return "", fmt.Errorf("Mongo password is undefined")
	}

	return fmt.Sprintf(
		"mongodb://%s:%s@%s:%s",
		username, password, host, port,
	), nil
}

func GetConnectionTimeout() (time.Duration, error) {
	connTimeoutStr, exists := os.LookupEnv("MONGO_CONNECTION_TIMEOUT")
	if !exists {
		return -1, fmt.Errorf("Mongo connection timeout is undefined")
	}
	connTimeout, err := strconv.Atoi(connTimeoutStr)
	if err != nil {
		return -1, err
	}
	return time.Duration(connTimeout), nil
}
