package db

import "time"

type Team struct {
	ID            string
	Name          string
	CompetitionID string `db:"competition_id"`

	CreatedAt int `db:"created_at"`
	UpdatedAt int `db:"updated_at"`
}

func (db *sqlImpl) GetTeam(id string) (team Team, err error) {
	err = db.db.Get(&team, "SELECT * FROM teams WHERE id=$1", id)
	return team, err
}

func (db *sqlImpl) InsertTeam(team Team) (err error) {
	team.CreatedAt = int(time.Now().Unix())
	team.UpdatedAt = team.CreatedAt
	_, err = db.db.NamedExec(
		`INSERT INTO teams (id, name, competition_id, created_at, updated_at) VALUES (:id, :name, :competition_id, :created_at, :updated_at)`,
		team)
	return err
}

func (db *sqlImpl) GetTeamsForCompetition(competitionId string) (teams []Team, err error) {
	err = db.db.Select(&teams, "SELECT * FROM teams WHERE competition_id=$1 ORDER BY name ASC", competitionId)
	return teams, err
}

func (db *sqlImpl) UpdateTeam(team Team) error {
	_, err := db.db.NamedExec(
		"UPDATE teams SET name=:name, competition_id=:competition_id, updated_at=:updated_at WHERE id=:id",
		team)
	return err
}

func (db *sqlImpl) DeleteTeam(id string) error {
	_, err := db.db.Exec("DELETE FROM teams WHERE id=$1", id)
	return err
}
