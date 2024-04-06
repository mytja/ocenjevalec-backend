package db

import (
	"time"
)

type Competition struct {
	ID          string
	Name        string
	Status      int
	StartTime   int `db:"start_time"`
	Penalty     int `db:"penalty"`      // time penalty
	PenaltyEach int `db:"penalty_each"` // per how many minutes a penalty should be given

	CreatedAt int `db:"created_at"`
	UpdatedAt int `db:"updated_at"`
}

func (db *sqlImpl) GetCompetition(id string) (competition Competition, err error) {
	err = db.db.Get(&competition, "SELECT * FROM competitions WHERE id=$1", id)
	return competition, err
}

func (db *sqlImpl) InsertCompetition(competition Competition) (err error) {
	competition.CreatedAt = int(time.Now().Unix())
	competition.UpdatedAt = competition.CreatedAt
	_, err = db.db.NamedExec(
		`INSERT INTO competitions (id, name, status, start_time, penalty, penalty_each, created_at, updated_at) VALUES (:id, :name, :status, :start_time, :penalty, :penalty_each, :created_at, :updated_at)`,
		competition)
	return err
}

func (db *sqlImpl) GetCompetitions() (competitions []Competition, err error) {
	err = db.db.Select(&competitions, "SELECT * FROM competitions ORDER BY status ASC")
	return competitions, err
}

func (db *sqlImpl) UpdateCompetition(competition Competition) error {
	competition.UpdatedAt = int(time.Now().Unix())
	_, err := db.db.NamedExec(
		"UPDATE competitions SET name=:name, status=:status, start_time=:start_time, penalty=:penalty, penalty_each=:penalty_each, updated_at=:updated_at WHERE id=:id",
		competition)
	return err
}

func (db *sqlImpl) DeleteCompetition(id string) error {
	_, err := db.db.Exec("DELETE FROM competitions WHERE id=$1", id)
	return err
}
