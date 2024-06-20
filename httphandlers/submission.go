package httphandlers

import (
	"HTTP-boilerplate/ast"
	"HTTP-boilerplate/db"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func (server *httpImpl) GetSubmission(w http.ResponseWriter, r *http.Request) {
	user, err := server.db.CheckToken(GetToken(r))
	if err != nil {
		WriteForbiddenJWT(w)
		return
	}

	if !user.IsAdmin {
		WriteJSON(w, Response{Error: "Forbidden"}, http.StatusForbidden)
		return
	}

	submissionId := mux.Vars(r)["submission_id"]
	submission, err := server.db.GetSubmission(submissionId)
	if err != nil {
		WriteJSON(w, Response{Error: "Server error whilst fetching competition"}, http.StatusInternalServerError)
		return
	}

	WriteJSON(w, Response{Data: submission}, http.StatusOK)
}

func (server *httpImpl) NewSubmission(w http.ResponseWriter, r *http.Request) {
	user, err := server.db.CheckToken(GetToken(r))
	if err != nil {
		WriteForbiddenJWT(w)
		return
	}

	if !user.IsAdmin {
		WriteJSON(w, Response{Error: "Forbidden"}, http.StatusForbidden)
		return
	}

	submittedAfter, err := strconv.Atoi(r.FormValue("submitted_after"))
	if err != nil {
		WriteJSON(w, Response{Error: "Submitted_after is invalid"}, http.StatusBadRequest)
		return
	}

	problemId := mux.Vars(r)["problem_id"]
	problem, err := server.db.GetProblem(problemId)
	if err != nil {
		WriteJSON(w, Response{Error: "Server error whilst fetching problem"}, http.StatusInternalServerError)
		return
	}

	competition, err := server.db.GetCompetition(problem.CompetitionID)
	if err != nil {
		WriteJSON(w, Response{Error: "Server error whilst fetching competition"}, http.StatusInternalServerError)
		return
	}

	teamId := r.FormValue("team_id")
	team, err := server.db.GetTeam(teamId)
	if err != nil {
		WriteJSON(w, Response{Error: "Server error whilst fetching team"}, http.StatusInternalServerError)
		return
	}

	problems, err := server.db.GetPastSubmissionsForProblem(submittedAfter, team.ID, problem.ID)
	if err != nil {
		WriteJSON(w, Response{Error: "Server error whilst fetching problems"}, http.StatusInternalServerError)
		return
	}

	previousSubmission := r.FormValue("previous_submission_id")

	id := uuid.NewString()

	submissionS := r.FormValue("submission")
	submissionS = ast.MinifyString(submissionS)

	submission := db.Submission{
		ID:             id,
		Solution:       submissionS,
		Verdict:        "",
		Score:          0,
		SubmittedAfter: submittedAfter,
		SubmissionLog:  "",
		CompetitionID:  problem.CompetitionID,
		ProblemID:      problemId,
		TeamID:         teamId,
		Public:         false,
	}

	if previousSubmission == "" {
		marshal, err := json.Marshal(WSNewSubmission{
			MessageType: 0,
			Submission:  id,
			TeamName:    team.Name,
			TeamID:      team.ID,
			ProblemID:   problem.ID,
			ProblemName: problem.Name,
			MaxScore:    problem.Points,
		})
		if err == nil {
			server.hub.broadcast <- marshal
		}
	} else {
		marshal, err := json.Marshal(WSChangeSubmissionID{
			MessageType:   2,
			OldSubmission: previousSubmission,
			NewSubmission: id,
		})
		if err == nil {
			server.hub.broadcast <- marshal
		}
	}

	sub, err := ast.BuildAST(submissionS)
	if err != nil {
		submission.SubmissionLog = err.Error()
		submission.Verdict = "CF" // Compilation failure
		err = server.db.InsertSubmission(submission)
		if err != nil {
			WriteJSON(w, Response{Error: "Server error whilst inserting submission"}, http.StatusInternalServerError)
			return
		}
		WriteJSON(w, Response{Data: submission}, http.StatusCreated)
		return
	}
	subL := ast.ASTLength(sub, 0)

	sol, err := ast.BuildAST(problem.Solution)
	if err != nil {
		submission.SubmissionLog = err.Error()
		submission.Verdict = "SOL_CF" // Solution compilation failure
		err = server.db.InsertSubmission(submission)
		if err != nil {
			WriteJSON(w, Response{Error: "Server error whilst inserting submission"}, http.StatusInternalServerError)
			return
		}
		WriteJSON(w, Response{Data: submission}, http.StatusCreated)
		return
	}
	solL := ast.ASTLength(sol, 0)

	test, err := ast.TestDigitalSolution(sub, sol)
	if err != nil {
		fmt.Println("Digital solution testing failed", err.Error())
		return
	}

	solH := ast.HashAST(sol)
	subH := ast.HashAST(sub)

	points := 0
	if test.Verdict == "AC" {
		if solH == subH || subL <= solL {
			points = problem.Points
			test.EvaluationLog += fmt.Sprintf("Equality check passed! Contestant: %s (%d, len: %d), Judge: %s (%d, len: %d).\n", submissionS, subH, subL, problem.Solution, solH, solL)
		} else {
			points = int(float64(problem.Points) * 0.9)
			test.EvaluationLog += fmt.Sprintf("Applying PART verdict! Contestant: %s (%d, len: %d), Judge: %s (%d, len: %d).\n", submissionS, subH, subL, problem.Solution, solH, solL)
			test.Verdict = "PART"
		}
	} else if test.Verdict == "WA" {
		// 0.72 izhaja iz tega, da so točke deljene tako:
		// 1. del (90 % vseh točk):
		//    - 70 % 1. dela (skupaj 63 %) gre testnim primerom
		//    - preostalih 30 % točk prvega dela (skupaj 27 %) gre ekvivalenci vsem testnim primerom, ki v takem primeru ni zadoščena
		// 2. del (10 % vseh točk):
		//    - vseh 10 % gre temu, da imajo tekmovalci rešitev identično uradni
		points = int(float64(problem.Points) * 0.63 * (float64(test.CorrectTestCases) / float64(test.WrongTestCases+test.CorrectTestCases)))
		test.EvaluationLog += fmt.Sprintf("Wrong answer! Wrong test cases: %d, Correct test cases: %d. Applying partial points: %d!\n", test.WrongTestCases, test.CorrectTestCases, points)
	}

	newPoints := max(0, points-len(problems)*30)
	test.EvaluationLog += fmt.Sprintf("Applying penalty of %d points due to previous submissions! Points before: %d, Points after: %d.\n", len(problems)*30, points, newPoints)
	points = newPoints

	newPoints = max(0, points-(submittedAfter*competition.Penalty))
	test.EvaluationLog += fmt.Sprintf("Applying penalty of %d points due to time! Points before: %d, Points after: %d.\n", submittedAfter*competition.Penalty, points, newPoints)
	points = newPoints

	submission.Verdict = test.Verdict
	submission.SubmissionLog = test.EvaluationLog
	submission.Score = points

	err = server.db.InsertSubmission(submission)
	if err != nil {
		WriteJSON(w, Response{Error: "Server error whilst inserting submission"}, http.StatusInternalServerError)
		return
	}

	server.db.DeleteSubmission(previousSubmission)

	WriteJSON(w, Response{Data: submission}, http.StatusCreated)
}

func (server *httpImpl) UpdateSubmission(w http.ResponseWriter, r *http.Request) {
	user, err := server.db.CheckToken(GetToken(r))
	if err != nil {
		WriteForbiddenJWT(w)
		return
	}

	if !user.IsAdmin {
		WriteJSON(w, Response{Error: "Forbidden"}, http.StatusForbidden)
		return
	}

	submissionId := mux.Vars(r)["submission_id"]
	submission, err := server.db.GetSubmission(submissionId)
	if err != nil {
		WriteJSON(w, Response{Error: "Server error whilst fetching a submission"}, http.StatusInternalServerError)
		return
	}

	//posodobi := !submission.Public

	public, err := strconv.ParseBool(r.FormValue("public"))
	if err == nil {
		submission.Public = public
	}

	score, err := strconv.Atoi(r.FormValue("score"))
	if err == nil {
		if score < 0 {
			WriteJSON(w, Response{Error: "Score is invalid. Expected a non-negative number."}, http.StatusBadRequest)
			return
		}
		submission.Score = score
		submission.Verdict = "MAN"
		submission.SubmissionLog = "Submission was judged manually!"
	}

	team, err := server.db.GetTeam(submission.TeamID)
	if err != nil {
		WriteJSON(w, Response{Error: "Server error whilst fetching the team"}, http.StatusInternalServerError)
		return
	}

	problem, err := server.db.GetProblem(submission.ProblemID)
	if err != nil {
		WriteJSON(w, Response{Error: "Server error whilst fetching the problem"}, http.StatusInternalServerError)
		return
	}

	err = server.db.UpdateSubmission(submission)
	if err != nil {
		WriteJSON(w, Response{Error: "Server error whilst updating submission"}, http.StatusInternalServerError)
		return
	}

	submissions1, err := server.db.GetTeamSubmissionsForProblem(team.ID, problem.ID)
	if err != nil {
		WriteJSON(w, Response{Error: "Server error whilst fetching team submissions"}, http.StatusInternalServerError)
		return
	}

	if len(submissions1) == 0 {
		WriteJSON(w, Response{Error: "Server error whilst saving submission"}, http.StatusInternalServerError)
		return
	}

	// posodobljen submission ni zadnji! Ne pošlji sporočila na klient
	if submissions1[len(submissions1)-1].ID != submissionId {
		WriteJSON(w, Response{Data: "OK"}, http.StatusOK)
		return
	}

	past, err := server.db.GetPastSubmissionsForProblem(submission.SubmittedAfter, team.ID, problem.ID)
	if err != nil {
		WriteJSON(w, Response{Error: "Server error whilst fetching past submissions"}, http.StatusInternalServerError)
		return
	}

	totalScore := 0

	problems, err := server.db.GetProblemsForCompetition(submission.CompetitionID)
	if err != nil {
		WriteJSON(w, Response{Error: "Server error whilst fetching problems"}, http.StatusInternalServerError)
		return
	}

	for _, v := range problems {
		submissions, err := server.db.GetTeamSubmissionsForProblem(team.ID, v.ID)
		if err != nil {
			continue
		}
		if len(submissions) == 0 {
			continue
		}
		totalScore += submissions[len(submissions)-1].Score
	}

	// v teoriji lahko posodobimo tudi ekipo, čeprav niti ne
	// "futureproofing"
	marshal, err := json.Marshal(WSUpdateSubmissionStatus{
		MessageType:       1,
		Submission:        submissionId,
		TeamName:          team.Name,
		TeamID:            team.ID,
		ProblemID:         problem.ID,
		ProblemName:       problem.Name,
		Verdict:           submission.Verdict,
		Score:             submission.Score,
		MaxScore:          problem.Points,
		TotalScore:        totalScore,
		SubmissionsBefore: len(past),
	})
	if err == nil {
		server.hub.broadcast <- marshal
	}

	WriteJSON(w, Response{Data: "OK"}, http.StatusOK)
}

func (server *httpImpl) DeleteSubmission(w http.ResponseWriter, r *http.Request) {
	user, err := server.db.CheckToken(GetToken(r))
	if err != nil {
		WriteForbiddenJWT(w)
		return
	}

	if !user.IsAdmin {
		WriteJSON(w, Response{Error: "Forbidden"}, http.StatusForbidden)
		return
	}

	submissionId := mux.Vars(r)["submission_id"]

	err = server.db.DeleteSubmission(submissionId)
	if err != nil {
		WriteJSON(w, Response{Error: "Server error whilst deleting a submission"}, http.StatusInternalServerError)
		return
	}

	WriteJSON(w, Response{Data: "OK"}, http.StatusOK)
}
