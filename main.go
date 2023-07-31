package main

import (
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/mhdianrush/go-login-signup-auth/config"
	"github.com/mhdianrush/go-login-signup-auth/controllers"
	"github.com/sirupsen/logrus"
)

func main() {
	config.ConnectDB()

	routes := mux.NewRouter()

	routes.HandleFunc("/", controllers.Index)
	routes.HandleFunc("/login", controllers.Login)
	routes.HandleFunc("/logout", controllers.Logout)
	routes.HandleFunc("/register", controllers.Register)

	logger := logrus.New()

	file, err := os.OpenFile("application.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	logger.SetOutput(file)

	if err := godotenv.Load(); err != nil {
		logger.Printf("failed load env file %s", err.Error())
	}

	server := http.Server{
		Addr:    ":" + os.Getenv("SERVER_PORT"),
		Handler: routes,
	}
	if err = server.ListenAndServe(); err != nil {
		logger.Printf("failed connect to server %s", err.Error())
	}
	
	logger.Printf("server running on port %s", os.Getenv("SERVER_PORT"))
}
