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
}

type HTTP interface {
	// user.go
	Login(w http.ResponseWriter, r *http.Request)
	NewUser(w http.ResponseWriter, r *http.Request)
	ChangePassword(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
}

func NewHTTPInterface(logger *zap.SugaredLogger, db db.SQL, config db.Config) HTTP {
	return &httpImpl{
		logger: logger,
		db:     db,
		config: config,
	}
}
