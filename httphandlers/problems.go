package httphandlers

import (
	"HTTP-boilerplate/ast"
	"HTTP-boilerplate/db"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func (server *httpImpl) GetProblems(w http.ResponseWriter, r *http.Request) {
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

	problems, err := server.db.GetProblemsForCompetition(competition.ID)
	if err != nil {
		WriteJSON(w, Response{Data: "Error whilst fetching problems"}, http.StatusInternalServerError)
		return
	}

	WriteJSON(w, Response{Data: problems}, http.StatusOK)
}

func (server *httpImpl) NewProblem(w http.ResponseWriter, r *http.Request) {
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

	points, err := strconv.Atoi(r.FormValue("points"))
	if err != nil {
		WriteJSON(w, Response{Error: "Invalid points"}, http.StatusBadRequest)
		return
	}

	if points < 0 {
		WriteJSON(w, Response{Error: "Invalid points"}, http.StatusBadRequest)
		return
	}

	competitionId := mux.Vars(r)["competition_id"]
	competition, err := server.db.GetCompetition(competitionId)
	if err != nil {
		WriteJSON(w, Response{Error: "Server error whilst fetching competition"}, http.StatusInternalServerError)
		return
	}

	problems, err := server.db.GetProblemsForCompetition(competition.ID)
	if err != nil {
		WriteJSON(w, Response{Error: "Server error whilst fetching problems"}, http.StatusInternalServerError)
		return
	}

	problemPos := 0
	if len(problems) != 0 {
		problemPos = problems[len(problems)-1].Position + 1
	}

	solution := r.FormValue("solution")
	solution = ast.MinifyString(solution)
	_, err = ast.BuildAST(solution)
	if err != nil {
		WriteJSON(w, Response{Error: "AST build failed for the solution"}, http.StatusInternalServerError)
		return
	}

	id := uuid.NewString()

	problem := db.Problem{
		ID:            id,
		Name:          name,
		Solution:      solution,
		Position:      problemPos,
		Points:        points,
		CompetitionID: competition.ID,
		AuthorID:      user.ID,
	}

	err = server.db.InsertProblem(problem)
	if err != nil {
		WriteJSON(w, Response{Error: "Server error whilst inserting problem"}, http.StatusInternalServerError)
		return
	}

	WriteJSON(w, Response{Data: id}, http.StatusCreated)
}

func (server *httpImpl) UpdateProblem(w http.ResponseWriter, r *http.Request) {
	user, err := server.db.CheckToken(GetToken(r))
	if err != nil {
		WriteForbiddenJWT(w)
		return
	}

	if !user.IsAdmin {
		WriteJSON(w, Response{Error: "Forbidden"}, http.StatusForbidden)
		return
	}

	problemId := mux.Vars(r)["problem_id"]
	problem, err := server.db.GetProblem(problemId)
	if err != nil {
		WriteJSON(w, Response{Error: "Server error whilst fetching a problem"}, http.StatusInternalServerError)
		return
	}

	problems, err := server.db.GetProblemsForCompetition(problem.CompetitionID)
	if err != nil {
		WriteJSON(w, Response{Error: "Server error whilst fetching problems"}, http.StatusInternalServerError)
		return
	}

	// sej se ne more zgoditi
	if len(problems) == 0 {
		WriteJSON(w, Response{Error: "No problems are applicable to the competition"}, http.StatusInternalServerError)
		return
	}

	name := r.FormValue("name")
	if name != "" {
		problem.Name = name
	}

	solution := r.FormValue("solution")
	if solution != "" {
		solution = ast.MinifyString(solution)
		_, err = ast.BuildAST(solution)
		if err != nil {
			WriteJSON(w, Response{Error: "AST build failed for the solution"}, http.StatusInternalServerError)
			return
		}
		problem.Solution = solution
	}

	points, err := strconv.Atoi(r.FormValue("points"))
	if err == nil {
		if points < 0 {
			WriteJSON(w, Response{Error: "Points are invalid. Expected a non-negative number."}, http.StatusBadRequest)
			return
		}
		problem.Points = points
	}

	// to naj bo na koncu, saj posodabljamo druge probleme
	position, err := strconv.Atoi(r.FormValue("position"))
	if err == nil {
		if position >= len(problems) {
			position = len(problems) - 1
		}
		problems[position].Position = problem.Position
		err = server.db.UpdateProblem(problems[position])
		if err != nil {
			WriteJSON(w, Response{Error: "Server error whilst updating applied problem"}, http.StatusInternalServerError)
			return
		}
		problem.Position = position
	}

	err = server.db.UpdateProblem(problem)
	if err != nil {
		WriteJSON(w, Response{Error: "Server error whilst updating problem"}, http.StatusInternalServerError)
		return
	}

	WriteJSON(w, Response{Data: "OK"}, http.StatusOK)
}

func (server *httpImpl) DeleteProblem(w http.ResponseWriter, r *http.Request) {
	user, err := server.db.CheckToken(GetToken(r))
	if err != nil {
		WriteForbiddenJWT(w)
		return
	}

	if !user.IsAdmin {
		WriteJSON(w, Response{Error: "Forbidden"}, http.StatusForbidden)
		return
	}

	problemId := mux.Vars(r)["problem_id"]
	problem, err := server.db.GetProblem(problemId)
	if err != nil {
		WriteJSON(w, Response{Error: "Server error whilst fetching a problem"}, http.StatusInternalServerError)
		return
	}

	problems, err := server.db.GetProblemsForCompetition(problem.CompetitionID)
	if err != nil {
		WriteJSON(w, Response{Error: "Server error whilst fetching problems"}, http.StatusInternalServerError)
		return
	}

	for _, v := range problems {
		if v.Position <= problem.Position {
			continue
		}
		v.Position--
		server.db.UpdateProblem(v)
	}

	err = server.db.DeleteProblem(problemId)
	if err != nil {
		WriteJSON(w, Response{Error: "Server error whilst deleting a problem"}, http.StatusInternalServerError)
		return
	}

	WriteJSON(w, Response{Data: "OK"}, http.StatusOK)
}
