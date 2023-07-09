package main

import (
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/mhdianrush/go-login-signup-auth/config"
	"github.com/sirupsen/logrus"
)

func main() {
	config.ConnectDB()

	r := mux.NewRouter()

	logger := logrus.New()

	file, err := os.OpenFile("application.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}

	logger.SetOutput(file)

	logger.Println("Server Running on Port 8080")

	err = http.ListenAndServe(":8080", r)
	if err != nil {
		panic(err)
	}
}
