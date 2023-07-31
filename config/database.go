package config

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

var logger = logrus.New()

func ConnectDB() (*sql.DB, error) {
	if err := godotenv.Load(); err != nil {
		logger.Printf("failed load env file %s", err.Error())
	}

	dataSourceName := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s",
		os.Getenv("DATABASE_USER"), os.Getenv("DATABASE_PASSWORD"), os.Getenv("DATABASE_HOST"), os.Getenv("DATABASE_PORT"), os.Getenv("DATABASE_NAME"),
	)
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		logger.Printf("failed connect to database %s", err.Error())
	}

	logger.Println("database connected")

	return db, err
}
