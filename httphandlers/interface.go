package httphandlers

import (
	"HTTP-boilerplate/db"
	"go.uber.org/zap"
	"net/http"
)

type Response struct {
	Error   any  `json:"error"`
	Success bool `json:"success"`
	Data    any  `json:"data"`
}

type httpImpl struct {
	logger *zap.SugaredLogger
	db     db.SQL
	config db.Config
	hub    *Hub
}

type HTTP interface {
	// user.go
	Login(w http.ResponseWriter, r *http.Request)
	NewUser(w http.ResponseWriter, r *http.Request)
	ChangePassword(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)

	// submission.go
	GetSubmission(w http.ResponseWriter, r *http.Request)
	NewSubmission(w http.ResponseWriter, r *http.Request)
	UpdateSubmission(w http.ResponseWriter, r *http.Request)
	DeleteSubmission(w http.ResponseWriter, r *http.Request)

	// competitions.go
	GetCompetitions(w http.ResponseWriter, r *http.Request)
	NewCompetition(w http.ResponseWriter, r *http.Request)
	UpdateCompetition(w http.ResponseWriter, r *http.Request)
	DeleteCompetition(w http.ResponseWriter, r *http.Request)
	BuildCompetitionLeaderboard(w http.ResponseWriter, r *http.Request)

	// problems.go
	GetProblems(w http.ResponseWriter, r *http.Request)
	NewProblem(w http.ResponseWriter, r *http.Request)
	UpdateProblem(w http.ResponseWriter, r *http.Request)
	DeleteProblem(w http.ResponseWriter, r *http.Request)

	// teams.go
	GetTeams(w http.ResponseWriter, r *http.Request)
	NewTeam(w http.ResponseWriter, r *http.Request)
	UpdateTeam(w http.ResponseWriter, r *http.Request)
	DeleteTeam(w http.ResponseWriter, r *http.Request)

	// ws-client.go
	UpgradeConnection(w http.ResponseWriter, r *http.Request)
}

func NewHTTPInterface(logger *zap.SugaredLogger, db db.SQL, config db.Config, hub *Hub) HTTP {
	return &httpImpl{
		logger: logger,
		db:     db,
		config: config,
		hub:    hub,
	}
}
