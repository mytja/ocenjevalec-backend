package db

import (
	"time"
)

type Submission struct {
	ID             string
	Solution       string // submitted solution
	Verdict        string
	Score          int
	SubmittedAfter int    `db:"submitted_after"` // after how many minutes has it been submitted
	SubmissionLog  string `db:"submission_log"`
	CompetitionID  string `db:"competition_id"`
	ProblemID      string `db:"problem_id"`
	TeamID         string `db:"team_id"`
	Public         bool   // whether the submission is displayed on leaderboards

	CreatedAt int `db:"created_at"`
	UpdatedAt int `db:"updated_at"`
}

func (db *sqlImpl) GetSubmission(id string) (submission Submission, err error) {
	err = db.db.Get(&submission, "SELECT * FROM submissions WHERE id=$1", id)
	return submission, err
}

func (db *sqlImpl) InsertSubmission(submission Submission) (err error) {
	submission.CreatedAt = int(time.Now().Unix())
	submission.UpdatedAt = submission.CreatedAt
	_, err = db.db.NamedExec(
		`INSERT INTO submissions (id, solution, verdict, score, submitted_after, submission_log, competition_id, problem_id, team_id, public, created_at, updated_at) VALUES
(:id, :solution, :verdict, :score, :submitted_after, :submission_log, :competition_id, :problem_id, :team_id, :public, :created_at, :updated_at)`,
		submission)
	return err
}

func (db *sqlImpl) GetPastSubmissionsForProblem(submittedAfter int, teamId string, problemId string) (submissions []Submission, err error) {
	err = db.db.Select(&submissions, "SELECT * FROM submissions WHERE submitted_after < $1 AND team_id=$2 AND problem_id=$3 AND public=true ORDER BY submitted_after ASC", submittedAfter, teamId, problemId)
	return submissions, err
}

func (db *sqlImpl) GetTeamSubmissionsForProblem(teamId string, problemId string) (submissions []Submission, err error) {
	err = db.db.Select(&submissions, "SELECT * FROM submissions WHERE team_id=$1 AND problem_id=$2 AND public=true ORDER BY submitted_after ASC", teamId, problemId)
	return submissions, err
}

func (db *sqlImpl) UpdateSubmission(submission Submission) error {
	submission.UpdatedAt = int(time.Now().Unix())
	_, err := db.db.NamedExec(
		"UPDATE submissions SET solution=:solution, verdict=:verdict, score=:score, submitted_after=:submitted_after, submission_log=:submission_log, competition_id=:competition_id, problem_id=:problem_id, team_id=:team_id, public=:public, updated_at=:updated_at WHERE id=:id",
		submission)
	return err
}

func (db *sqlImpl) DeleteSubmission(id string) error {
	_, err := db.db.Exec("DELETE FROM submissions WHERE id=$1", id)
	return err
}
