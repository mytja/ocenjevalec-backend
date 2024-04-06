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

	hub := httphandlers.NewHub()
	go hub.Run()

	httphandler := httphandlers.NewHTTPInterface(sugared, database, config, hub)

	sugared.Info("Database created successfully")

	r := mux.NewRouter()
	r.HandleFunc("/user/new", httphandler.NewUser).Methods("POST")
	r.HandleFunc("/user/login", httphandler.Login).Methods("POST")
	r.HandleFunc("/user/get/password_change", httphandler.ChangePassword).Methods("PATCH")
	r.HandleFunc("/user/logout", httphandler.Logout).Methods("POST")

	r.HandleFunc("/submission/{submission_id}", httphandler.GetSubmission).Methods("GET")
	r.HandleFunc("/problem/{problem_id}/submission", httphandler.NewSubmission).Methods("POST")
	r.HandleFunc("/submission/{submission_id}", httphandler.UpdateSubmission).Methods("PATCH")
	r.HandleFunc("/submission/{submission_id}", httphandler.DeleteSubmission).Methods("DELETE")

	r.HandleFunc("/competitions", httphandler.GetCompetitions).Methods("GET")
	r.HandleFunc("/competitions", httphandler.NewCompetition).Methods("POST")
	r.HandleFunc("/competition/{competition_id}", httphandler.UpdateCompetition).Methods("PATCH")
	r.HandleFunc("/competition/{competition_id}", httphandler.DeleteCompetition).Methods("DELETE")
	r.HandleFunc("/competition/{competition_id}", httphandler.BuildCompetitionLeaderboard).Methods("GET")

	r.HandleFunc("/competition/{competition_id}/websocket", httphandler.UpgradeConnection).Methods("GET")

	r.HandleFunc("/competition/{competition_id}/problems", httphandler.GetProblems).Methods("GET")
	r.HandleFunc("/competition/{competition_id}/problems", httphandler.NewProblem).Methods("POST")
	r.HandleFunc("/problem/{problem_id}", httphandler.UpdateProblem).Methods("PATCH")
	r.HandleFunc("/problem/{problem_id}", httphandler.DeleteProblem).Methods("DELETE")

	r.HandleFunc("/competition/{competition_id}/teams", httphandler.GetTeams).Methods("GET")
	r.HandleFunc("/competition/{competition_id}/teams", httphandler.NewTeam).Methods("POST")
	r.HandleFunc("/team/{team_id}", httphandler.UpdateTeam).Methods("PATCH")
	r.HandleFunc("/team/{team_id}", httphandler.DeleteTeam).Methods("DELETE")

	o := cors.Options{
		AllowedMethods:   []string{"POST", "GET", "DELETE", "PATCH", "PUT"},
		AllowCredentials: true,
	}

	if config.Debug {
		o.AllowedOrigins = []string{"http://localhost:5173", "http://127.0.0.1:5173"}
	} else {
		o.AllowedOrigins = []string{"https://ocenjevalec.beziapp.si"}
	}

	c := cors.New(o)

	sugared.Infof("Serving at %s", config.Host)

	err = http.ListenAndServe(config.Host, c.Handler(r))
	if err != nil {
		sugared.Fatal(err.Error())
	}

	sugared.Info("Done serving...")
}
