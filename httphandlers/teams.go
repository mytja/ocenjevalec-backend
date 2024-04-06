package httphandlers

import (
	"HTTP-boilerplate/db"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
)

func (server *httpImpl) GetTeams(w http.ResponseWriter, r *http.Request) {
	user, err := server.db.CheckToken(GetToken(r))
	if err != nil {
		WriteForbiddenJWT(w)
		return
	}

	if !user.IsAdmin {
		WriteJSON(w, Response{Error: "Forbidden"}, http.StatusForbidden)
		return
	}

	competitionId := mux.Vars(r)["competition_id"]
	competition, err := server.db.GetCompetition(competitionId)
	if err != nil {
		WriteJSON(w, Response{Error: "Server error whilst fetching competition"}, http.StatusInternalServerError)
		return
	}

	teams, err := server.db.GetTeamsForCompetition(competition.ID)
	if err != nil {
		WriteJSON(w, Response{Data: "Error whilst fetching problems"}, http.StatusInternalServerError)
		return
	}

	if teams == nil {
		teams = make([]db.Team, 0)
	}

	WriteJSON(w, Response{Data: teams}, http.StatusOK)
}

func (server *httpImpl) NewTeam(w http.ResponseWriter, r *http.Request) {
	user, err := server.db.CheckToken(GetToken(r))
	if err != nil {
		WriteForbiddenJWT(w)
		return
	}

	if !user.IsAdmin {
		WriteJSON(w, Response{Error: "Forbidden"}, http.StatusForbidden)
		return
	}

	name := r.FormValue("name")
	if name == "" {
		WriteJSON(w, Response{Error: "Invalid name"}, http.StatusBadRequest)
		return
	}

	competitionId := mux.Vars(r)["competition_id"]
	competition, err := server.db.GetCompetition(competitionId)
	if err != nil {
		WriteJSON(w, Response{Error: "Server error whilst fetching competition"}, http.StatusInternalServerError)
		return
	}

	id := uuid.NewString()

	team := db.Team{
		ID:            id,
		Name:          name,
		CompetitionID: competition.ID,
	}

	err = server.db.InsertTeam(team)
	if err != nil {
		WriteJSON(w, Response{Error: "Server error whilst inserting team"}, http.StatusInternalServerError)
		return
	}

	WriteJSON(w, Response{Data: id}, http.StatusCreated)
}

func (server *httpImpl) UpdateTeam(w http.ResponseWriter, r *http.Request) {
	user, err := server.db.CheckToken(GetToken(r))
	if err != nil {
		WriteForbiddenJWT(w)
		return
	}

	if !user.IsAdmin {
		WriteJSON(w, Response{Error: "Forbidden"}, http.StatusForbidden)
		return
	}

	teamId := mux.Vars(r)["team_id"]
	team, err := server.db.GetTeam(teamId)
	if err != nil {
		WriteJSON(w, Response{Error: "Server error whilst fetching a team"}, http.StatusInternalServerError)
		return
	}

	team.Name = r.FormValue("name")
	if team.Name == "" {
		WriteJSON(w, Response{Error: "Invalid name"}, http.StatusBadRequest)
		return
	}

	err = server.db.UpdateTeam(team)
	if err != nil {
		WriteJSON(w, Response{Error: "Server error whilst updating problem"}, http.StatusInternalServerError)
		return
	}

	WriteJSON(w, Response{Data: "OK"}, http.StatusOK)
}

func (server *httpImpl) DeleteTeam(w http.ResponseWriter, r *http.Request) {
	user, err := server.db.CheckToken(GetToken(r))
	if err != nil {
		WriteForbiddenJWT(w)
		return
	}

	if !user.IsAdmin {
		WriteJSON(w, Response{Error: "Forbidden"}, http.StatusForbidden)
		return
	}

	teamId := mux.Vars(r)["team_id"]

	err = server.db.DeleteTeam(teamId)
	if err != nil {
		WriteJSON(w, Response{Error: "Server error whilst deleting a team"}, http.StatusInternalServerError)
		return
	}

	WriteJSON(w, Response{Data: "OK"}, http.StatusOK)
}
