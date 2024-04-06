package db

const schema string = `
CREATE TABLE IF NOT EXISTS users (
	id                       VARCHAR(40)    PRIMARY KEY,
	username                 VARCHAR(250)   UNIQUE NOT NULL,
	pass                     VARCHAR(250)   NOT NULL,
	login_token              VARCHAR(400),
	is_admin                 BOOLEAN        NOT NULL,
	is_locked                BOOLEAN        NOT NULL,
	
	created_at               INTEGER,
	updated_at               INTEGER
);

CREATE TABLE IF NOT EXISTS competitions (
	id                       VARCHAR(40)    PRIMARY KEY,
	name                     VARCHAR(250)   NOT NULL,
	status                   INTEGER        NOT NULL,
	start_time               INTEGER        NOT NULL,
	penalty                  INTEGER        NOT NULL,
	penalty_each             INTEGER        NOT NULL,
	
	created_at               INTEGER,
	updated_at               INTEGER
);

CREATE TABLE IF NOT EXISTS problems (
	id                       VARCHAR(40)    PRIMARY KEY,
	name                     VARCHAR(250)   NOT NULL,
	solution                 VARCHAR(500)   NOT NULL,
	position                 INTEGER        NOT NULL,
	points                   INTEGER        NOT NULL,
	competition_id           VARCHAR(40)    NOT NULL,
    author_id                VARCHAR(40)    NOT NULL,
	
	created_at               INTEGER,
	updated_at               INTEGER
);

CREATE TABLE IF NOT EXISTS submissions (
	id                       VARCHAR(40)    PRIMARY KEY,
	solution                 VARCHAR(500)   NOT NULL,
	verdict                  VARCHAR(40)    NOT NULL,
	score                    INTEGER        NOT NULL,
    submitted_after          INTEGER        NOT NULL,
    submission_log           VARCHAR(20000) NOT NULL,
	competition_id           VARCHAR(40)    NOT NULL,
    problem_id               VARCHAR(40)    NOT NULL,
    team_id                  VARCHAR(40)    NOT NULL,
    public                   BOOLEAN        NOT NULL,
	
	created_at               INTEGER,
	updated_at               INTEGER
);

CREATE TABLE IF NOT EXISTS teams (
	id                       VARCHAR(40)    PRIMARY KEY,
	name                     VARCHAR(250)   NOT NULL,
	competition_id           VARCHAR(40)    NOT NULL,
	
	created_at               INTEGER,
	updated_at               INTEGER
);
`
