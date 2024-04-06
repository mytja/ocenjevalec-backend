package db

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
)

type sqlImpl struct {
	db     *sqlx.DB
	logger *zap.SugaredLogger
}

func (db *sqlImpl) Init() {
	db.db.MustExec(schema)
}

type SQL interface {
	CheckToken(loginToken string) (User, error)
	GetRandomToken(currentUser User) (string, error)

	Init()
	Exec(query string) error

	GetUser(id string) (user User, err error)
	GetUserByLoginToken(loginToken string) (user User, err error)
	InsertUser(user User) (err error)
	GetUserByUsername(username string) (user User, err error)
	CheckIfAdminIsCreated() bool
	GetUsers() (users []User, err error)
	UpdateUser(user User) error
	DeleteUser(ID string) error

	GetCompetition(id string) (competition Competition, err error)
	InsertCompetition(competition Competition) (err error)
	GetCompetitions() (competitions []Competition, err error)
	UpdateCompetition(competition Competition) error
	DeleteCompetition(id string) error

	GetProblem(id string) (problem Problem, err error)
	InsertProblem(problem Problem) (err error)
	GetProblemsForCompetition(competitionId string) (problems []Problem, err error)
	UpdateProblem(problem Problem) error
	DeleteProblem(id string) error

	GetSubmission(id string) (submission Submission, err error)
	InsertSubmission(submission Submission) (err error)
	GetPastSubmissionsForProblem(submittedAfter int, teamId string, problemId string) (submissions []Submission, err error)
	GetTeamSubmissionsForProblem(teamId string, problemId string) (submissions []Submission, err error)
	UpdateSubmission(submission Submission) error
	DeleteSubmission(id string) error

	GetTeam(id string) (team Team, err error)
	InsertTeam(team Team) (err error)
	GetTeamsForCompetition(competitionId string) (teams []Team, err error)
	UpdateTeam(team Team) error
	DeleteTeam(id string) error
}

func NewSQL(driver string, drivername string, logger *zap.SugaredLogger) (SQL, error) {
	db, err := sqlx.Connect(driver, drivername)
	return &sqlImpl{
		db:     db,
		logger: logger,
	}, err
}

func (db *sqlImpl) Exec(query string) error {
	_, err := db.db.Exec(query)
	return err
}
