package postgres

import (
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db   *gorm.DB
	once sync.Once
)

func GetDB() *gorm.DB {
	once.Do(func() {
		dsn, err := GetDSN()
		if err != nil {
			panic("Failed to get DSN Postgres: " + err.Error())
		}

		conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			panic("Failed to open connection to Postgres: " + err.Error())
		}

		db = conn
	})

	return db
}
