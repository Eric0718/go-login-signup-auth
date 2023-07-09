package config

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

func ConnectDB() (*sql.DB, error) {
	logger := logrus.New()

	db, err := sql.Open("mysql", "root:admin@tcp(127.0.0.1:3306)/go_login_signup_auth?parseTime=true")
	if err != nil {
		panic(err)
	}

	logger.Println("Database Connected")

	return db, err
}
