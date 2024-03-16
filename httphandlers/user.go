package httphandlers

import (
	"HTTP-boilerplate/db"
	"database/sql"
	"github.com/google/uuid"
	"net/http"
	"time"
)

type TokenResponse struct {
	UserID   string `json:"user_id"`
	Token    string `json:"token"`
	IsAdmin  bool   `json:"is_admin"`
	Username string `json:"username"`
}

func (server *httpImpl) Login(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	pass := r.FormValue("pass")
	// Check if password is valid
	user, err := server.db.GetUserByUsername(username)
	if err != nil {
		WriteJSON(w, Response{Data: "Failed while retrieving the user", Success: false, Error: err.Error()}, http.StatusInternalServerError)
		return
	}

	if user.IsLocked {
		WriteJSON(w, Response{Data: "You are LOCKED. This can be from several different reasons, such as abuse of your account. You cannot log in until administrators unlock your account.", Success: false}, http.StatusForbidden)
		return
	}

	hashCorrect := db.CheckHash(pass, user.Password)
	if !hashCorrect {
		WriteJSON(w, Response{Data: "Hashes don't match... Is the password correct?", Success: false}, http.StatusForbidden)
		return
	}

	var token string

	if user.LoginToken == "" {
		token, err = server.db.GetRandomToken(user)
		if err != nil {
			WriteJSON(w, Response{Error: err.Error(), Success: false}, http.StatusInternalServerError)
			return
		}
	} else {
		token = user.LoginToken
	}

	WriteJSON(w, Response{Data: TokenResponse{Token: token, UserID: user.ID, Username: user.Username, IsAdmin: user.IsAdmin}, Success: true}, http.StatusOK)
}

func (server *httpImpl) NewUser(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	pass := r.FormValue("pass")
	if username == "" || pass == "" {
		WriteJSON(w, Response{Data: "Bad Request. A parameter isn't provided.", Success: false}, http.StatusBadRequest)
		return
	}

	// Check if user is already existing inside the database
	var userCreated = true
	_, err := server.db.GetUserByUsername(username)
	if err != nil {
		if err == sql.ErrNoRows {
			userCreated = false
		} else {
			WriteJSON(w, Response{Error: err.Error(), Data: "Could not retrieve user from database", Success: false}, http.StatusInternalServerError)
			return
		}
	}

	if userCreated == true {
		WriteJSON(w, Response{Data: "User is already in database", Success: false}, http.StatusUnprocessableEntity)
		return
	}

	// It's important to hash the password before committing it to the database ;)
	password, err := db.HashPassword(pass)
	if err != nil {
		WriteJSON(w, Response{Error: err.Error(), Data: "Failed to hash your password", Success: false}, http.StatusInternalServerError)
		return
	}

	isAdmin := !server.db.CheckIfAdminIsCreated()

	currentTime := int(time.Now().Unix())

	user := db.User{
		ID:         uuid.NewString(),
		Username:   username,
		Password:   password,
		IsAdmin:    isAdmin,
		LoginToken: "",
		IsLocked:   false,
		CreatedAt:  currentTime,
		UpdatedAt:  currentTime,
	}

	err = server.db.InsertUser(user)
	if err != nil {
		WriteJSON(w, Response{Error: err.Error(), Data: "Failed to commit new user to database", Success: false}, http.StatusInternalServerError)
		return
	}

	WriteJSON(w, Response{Data: "Success", Success: true}, http.StatusCreated)
}

func (server *httpImpl) ChangePassword(w http.ResponseWriter, r *http.Request) {
	user, err := server.db.CheckToken(GetAuthorizationJWT(r))
	if err != nil {
		WriteForbiddenJWT(w)
		return
	}

	oldPass := r.FormValue("oldPassword")
	if !db.CheckHash(oldPass, user.Password) {
		WriteJSON(w, Response{Data: "Wrong password", Success: false}, http.StatusForbidden)
		return
	}

	password, err := db.HashPassword(r.FormValue("password"))
	if err != nil {
		return
	}

	user.Password = password
	user.UpdatedAt = int(time.Now().Unix())

	err = server.db.UpdateUser(user)
	if err != nil {
		WriteJSON(w, Response{Data: "Error while updating the database", Success: false, Error: err.Error()}, http.StatusInternalServerError)
		return
	}

	WriteJSON(w, Response{Data: "OK", Success: true}, http.StatusOK)
}

// Logout takes care of invalidating existing session tokens
func (server *httpImpl) Logout(w http.ResponseWriter, r *http.Request) {
	user, err := server.db.CheckToken(GetAuthorizationJWT(r))
	if err != nil {
		WriteForbiddenJWT(w)
		return
	}

	user.LoginToken = ""
	user.UpdatedAt = int(time.Now().Unix())

	err = server.db.UpdateUser(user)
	if err != nil {
		WriteJSON(w, Response{Data: "Error while updating the database", Success: false, Error: err.Error()}, http.StatusInternalServerError)
		return
	}

	WriteJSON(w, Response{Data: "OK", Success: true}, http.StatusOK)
}
