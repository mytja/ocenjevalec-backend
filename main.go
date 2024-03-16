package main

import (
	"HTTP-boilerplate/db"
	"HTTP-boilerplate/httphandlers"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"go.uber.org/zap"
	"net/http"
	"os"
)

func main() {
	fmt.Println("Starting server...")

	var logger *zap.Logger
	var err error

	config, err := db.GetConfig()
	if err != nil {
		panic("Error while retrieving config: " + err.Error())
		return
	}

	if config.Debug {
		logger, err = zap.NewDevelopment()
	} else {
		logger, err = zap.NewProduction()
	}

	if err != nil {
		panic(err.Error())
		return
	}

	sugared := logger.Sugar()

	if _, err := os.Stat("database"); os.IsNotExist(err) {
		err := os.Mkdir("database", os.ModePerm)
		if err != nil {
			panic("Cannot create a directory database")
		}
	}

	database, err := db.NewSQL(config.DatabaseName, config.DatabaseConfig, sugared)
	database.Init()

	if err != nil {
		sugared.Fatal("Error while creating database: ", err.Error())
		return
	}

	httphandler := httphandlers.NewHTTPInterface(sugared, database, config)

	sugared.Info("Database created successfully")

	r := mux.NewRouter()
	r.HandleFunc("/user/new", httphandler.NewUser).Methods("POST")
	r.HandleFunc("/user/login", httphandler.Login).Methods("POST")
	r.HandleFunc("/user/get/password_change", httphandler.ChangePassword).Methods("PATCH")
	r.HandleFunc("/user/logout", httphandler.Logout).Methods("POST")

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"}, // All origins
		AllowedHeaders: []string{"Authorization"},
		AllowedMethods: []string{"POST", "GET", "DELETE", "PATCH", "PUT"},
	})

	sugared.Infof("Serving at %s", config.Host)

	err = http.ListenAndServe(config.Host, c.Handler(r))
	if err != nil {
		sugared.Fatal(err.Error())
	}

	sugared.Info("Done serving...")
}
