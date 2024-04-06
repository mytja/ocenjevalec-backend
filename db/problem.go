package db

import "time"

type Problem struct {
	ID            string
	Name          string
	Solution      string
	Position      int
	Points        int
	CompetitionID string `db:"competition_id"`
	AuthorID      string `db:"author_id"`

	CreatedAt int `db:"created_at"`
	UpdatedAt int `db:"updated_at"`
}

func (db *sqlImpl) GetProblem(id string) (problem Problem, err error) {
	err = db.db.Get(&problem, "SELECT * FROM problems WHERE id=$1", id)
	return problem, err
}

func (db *sqlImpl) InsertProblem(problem Problem) (err error) {
	problem.CreatedAt = int(time.Now().Unix())
	problem.UpdatedAt = problem.CreatedAt
	_, err = db.db.NamedExec(
		`INSERT INTO problems (id, name, solution, position, points, competition_id, author_id, created_at, updated_at) VALUES (:id, :name, :solution, :position, :points, :competition_id, :author_id, :created_at, :updated_at)`,
		problem)
	return err
}

func (db *sqlImpl) GetProblemsForCompetition(competitionId string) (problems []Problem, err error) {
	err = db.db.Select(&problems, "SELECT * FROM problems WHERE competition_id=$1 ORDER BY position", competitionId)
	return problems, err
}

func (db *sqlImpl) UpdateProblem(problem Problem) error {
	_, err := db.db.NamedExec(
		"UPDATE problems SET name=:name, solution=:solution, position=:position, points=:points, updated_at=:updated_at, competition_id=:competition_id, author_id=:author_id WHERE id=:id",
		problem)
	return err
}

func (db *sqlImpl) DeleteProblem(id string) error {
	_, err := db.db.Exec("DELETE FROM problems WHERE id=$1", id)
	return err
}
