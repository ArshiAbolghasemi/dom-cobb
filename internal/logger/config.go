package logger

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

func GetCollection() (string, error) {
	collection, exists := os.LookupEnv("MONGO_LOG_COLLECTION")
	if !exists {
		return "", fmt.Errorf("Log collection mongo is undefined")
	}
	return collection, nil
}

func GetWriteTimeOut() (time.Duration, error) {
	timeoutStr, exists := os.LookupEnv("MONGO_LOG_WRITE_TIME_OUT")
	if !exists {
		return -1, fmt.Errorf("Log write timeout mongo is undefined")
	}
	timeout, err := strconv.Atoi(timeoutStr)
	if err != nil {
		return -1, err
	}
	return time.Duration(timeout), nil
}
