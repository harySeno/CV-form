package database

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var DB *gorm.DB

// InitDB to initiate database connection to postgres server
func InitDB() (*gorm.DB, error) {
	db, err := gorm.Open("postgres", "host=127.0.0.1 port=5432 user=testuser password=1234 dbname=testdb sslmode=disable")
	if err != nil {
		return nil, err
	}

	DB = db
	return db, nil
}
