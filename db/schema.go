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
`
