package store

import (
	"errors"
	"fmt"

	"github.com/shahin-bayat/go-scraper/util"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresStore() (*gorm.DB, error) {

	host := util.GetEnvVariable("DB_HOST")
	user := util.GetEnvVariable("DB_USER")
	password := util.GetEnvVariable("DB_PASSWORD")
	dbname := util.GetEnvVariable("DB_NAME")
	port := util.GetEnvVariable("DB_PORT")
	sslmode := util.GetEnvVariable("DB_SSLMODE")
	timezone := util.GetEnvVariable("DB_TIMEZONE")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s", host, user, password, dbname, port, sslmode, timezone)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, errors.New("failed to connect database")
	}

	return db, nil
}
