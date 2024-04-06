package httphandlers

import (
	"HTTP-boilerplate/db"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
	"sort"
	"strconv"
	"time"
)

func (server *httpImpl) GetCompetitions(w http.ResponseWriter, r *http.Request) {
	user, err := server.db.CheckToken(GetToken(r))
	if err != nil {
		WriteForbiddenJWT(w)
		return
	}

	if !user.IsAdmin {
		WriteJSON(w, Response{Error: "Forbidden"}, http.StatusForbidden)
		return
	}

	competitions, err := server.db.GetCompetitions()
	if err != nil {
		WriteJSON(w, Response{Data: "Error whilst fetching competitions"}, http.StatusInternalServerError)
		return
	}

	if competitions == nil {
		competitions = make([]db.Competition, 0)
	}

	WriteJSON(w, Response{Data: competitions}, http.StatusOK)
}

func (server *httpImpl) NewCompetition(w http.ResponseWriter, r *http.Request) {
	user, err := server.db.CheckToken(GetToken(r))
	if err != nil {
		WriteForbiddenJWT(w)
		return
	}

	if !user.IsAdmin {
		WriteJSON(w, Response{Error: "Forbidden"}, http.StatusForbidden)
		return
	}

	penalty, err := strconv.Atoi(r.FormValue("penalty"))
	if err != nil {
		WriteJSON(w, Response{Error: "Penalty is invalid. Expected float."}, http.StatusBadRequest)
		return
	}

	if penalty < 0 {
		WriteJSON(w, Response{Error: "Penalty is invalid. Expected non-negative number."}, http.StatusBadRequest)
		return
	}

	penaltyEach, err := strconv.Atoi(r.FormValue("penalty_each"))
	if err != nil {
		WriteJSON(w, Response{Error: "Penalty_each is invalid. Expected float."}, http.StatusBadRequest)
		return
	}

	if penaltyEach <= 0 {
		WriteJSON(w, Response{Error: "Penalty_each is invalid. Expected positive number."}, http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	if name == "" {
		WriteJSON(w, Response{Error: "Name is invalid. Expected non-empty string."}, http.StatusBadRequest)
		return
	}

	id := uuid.NewString()

	competition := db.Competition{
		ID:          id,
		Name:        name,
		Status:      0,
		Penalty:     penalty,
		PenaltyEach: penaltyEach,
	}

	err = server.db.InsertCompetition(competition)
	if err != nil {
		WriteJSON(w, Response{Error: "Server error whilst inserting competition", Data: err.Error()}, http.StatusInternalServerError)
		return
	}

	WriteJSON(w, Response{Data: id}, http.StatusCreated)
}

func (server *httpImpl) UpdateCompetition(w http.ResponseWriter, r *http.Request) {
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

	status, err := strconv.Atoi(r.FormValue("status"))
	if err == nil {
		if status < 0 || status > 2 {
			WriteJSON(w, Response{Error: "Status is invalid. Expected integer on interval [0, 2]."}, http.StatusBadRequest)
			return
		}
		if status == 1 {
			competition.StartTime = int(time.Now().Unix())
		}
		competition.Status = status
	}

	name := r.FormValue("name")
	if name != "" {
		competition.Name = name
	}

	penalty, err := strconv.Atoi(r.FormValue("penalty"))
	if err == nil {
		if penalty < 0 {
			WriteJSON(w, Response{Error: "Penalty is invalid. Expected a non-negative number."}, http.StatusBadRequest)
			return
		}
		competition.Penalty = penalty
	}

	penaltyEach, err := strconv.Atoi(r.FormValue("penalty_each"))
	if err == nil {
		if penaltyEach <= 0 {
			WriteJSON(w, Response{Error: "Penalty is invalid. Expected a positive number."}, http.StatusBadRequest)
			return
		}
		competition.PenaltyEach = penaltyEach
	}

	err = server.db.UpdateCompetition(competition)
	if err != nil {
		WriteJSON(w, Response{Error: "Server error whilst updating competition"}, http.StatusInternalServerError)
		return
	}

	WriteJSON(w, Response{Data: "OK"}, http.StatusOK)
}

func (server *httpImpl) DeleteCompetition(w http.ResponseWriter, r *http.Request) {
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

	err = server.db.DeleteCompetition(competitionId)
	if err != nil {
		WriteJSON(w, Response{Error: "Server error whilst deleting competition"}, http.StatusInternalServerError)
		return
	}

	WriteJSON(w, Response{Data: "OK"}, http.StatusOK)
}

type LeaderboardProblem struct {
	LatestSubmission  db.Submission
	SubmissionsBefore int
}

type LeaderboardTeam struct {
	Team       db.Team
	Problems   []*LeaderboardProblem
	TotalScore int
}

type Leaderboard struct {
	Competition db.Competition
	Problems    []db.Problem
	Teams       []LeaderboardTeam
}

func (server *httpImpl) BuildCompetitionLeaderboard(w http.ResponseWriter, r *http.Request) {
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

	problems, err := server.db.GetProblemsForCompetition(competitionId)
	if err != nil {
		return
	}

	if problems == nil {
		problems = make([]db.Problem, 0)
	}

	teams, err := server.db.GetTeamsForCompetition(competitionId)
	if err != nil {
		return
	}

	lbteams := make([]LeaderboardTeam, 0)

	for _, team := range teams {
		lbteam := LeaderboardTeam{
			Team:       team,
			Problems:   make([]*LeaderboardProblem, len(problems)),
			TotalScore: 0,
		}
		for l, problem := range problems {
			submissions, err := server.db.GetTeamSubmissionsForProblem(team.ID, problem.ID)
			if err != nil {
				continue
			}
			if submissions == nil || len(submissions) == 0 {
				lbteam.Problems[l] = nil
				continue
			}
			lbteam.Problems[l] = &LeaderboardProblem{}
			problems[l].Solution = ""
			lbteam.Problems[l].LatestSubmission = submissions[len(submissions)-1]
			lbteam.Problems[l].SubmissionsBefore = len(submissions) - 1
			lbteam.TotalScore += submissions[len(submissions)-1].Score
		}
		lbteams = append(lbteams, lbteam)
	}

	sort.Slice(lbteams, func(i, j int) bool {
		return lbteams[i].TotalScore > lbteams[j].TotalScore
	})

	WriteJSON(w, Response{Data: Leaderboard{
		Competition: competition,
		Problems:    problems,
		Teams:       lbteams,
	}}, http.StatusOK)
}
