package db

type User struct {
	ID         string
	Username   string
	Password   string `db:"pass"`
	IsAdmin    bool   `db:"is_admin"`
	LoginToken string `db:"login_token"`
	IsLocked   bool   `db:"is_locked"`

	CreatedAt int `db:"created_at"`
	UpdatedAt int `db:"updated_at"`
}

func (db *sqlImpl) GetUser(id string) (user User, err error) {
	err = db.db.Get(&user, "SELECT * FROM users WHERE id=$1", id)
	return user, err
}

func (db *sqlImpl) GetUserByLoginToken(loginToken string) (user User, err error) {
	err = db.db.Get(&user, "SELECT * FROM users WHERE login_token=$1", loginToken)
	return user, err
}

func (db *sqlImpl) GetUserByUsername(email string) (user User, err error) {
	err = db.db.Get(&user, "SELECT * FROM users WHERE username=$1", email)
	return user, err
}

func (db *sqlImpl) InsertUser(user User) (err error) {
	_, err = db.db.NamedExec(
		`INSERT INTO users (id, username, pass, is_admin, login_token, is_locked, created_at, updated_at) 
VALUES (:id, :username, :pass, :is_admin, :login_token, :is_locked, :created_at, :updated_at)`,
		user)
	return err
}

func (db *sqlImpl) CheckIfAdminIsCreated() bool {
	var users []User
	err := db.db.Select(&users, "SELECT * FROM users")
	if err != nil {
		// Return true, as we don't want all the users, on some internal error, to become administrators
		return true
	}
	return len(users) > 0
}

func (db *sqlImpl) GetUsers() (users []User, err error) {
	err = db.db.Select(&users, "SELECT * FROM users ORDER BY id ASC")
	return users, err
}

func (db *sqlImpl) UpdateUser(user User) error {
	_, err := db.db.NamedExec(
		"UPDATE users SET pass=:pass, is_admin=:is_admin, login_token=:login_token, is_locked=:is_locked, created_at=:created_at, updated_at=:updated_at WHERE id=:id",
		user)
	return err
}

func (db *sqlImpl) DeleteUser(ID string) error {
	_, err := db.db.Exec("DELETE FROM users WHERE id=$1", ID)
	return err
}
