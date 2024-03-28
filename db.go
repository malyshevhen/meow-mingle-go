package main

import (
	"database/sql"
	"log"

	"github.com/go-sql-driver/mysql"
)

type MySQLStorage struct {
	db *sql.DB
}

func NewMySQLStorage(config mysql.Config) *MySQLStorage {
	db, err := sql.Open("mysql", config.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to the MySQL DB")

	return &MySQLStorage{db: db}
}

func (m *MySQLStorage) Init() (*sql.DB, error) {
	// Init the tables
	if err := m.createUsersTable(); err != nil {
		return nil, err
	}

	if err := m.createProjectsTable(); err != nil {
		return nil, err
	}

	if err := m.createTasksTable(); err != nil {
		return nil, err
	}

	return m.db, nil
}

func (m *MySQLStorage) createUsersTable() error {
	_, err := m.db.Exec(`
CREATE TABLE IF NOT EXISTS users (
	id 					INT UNSIGNED NOT NULL AUTO_INCREMENT,
	email 			VARCHAR(255) NOT NULL,
	first_name 	VARCHAR(255) NOT NULL,
	last_name 	VARCHAR(255) NOT NULL,
	password 		VARCHAR(255) NOT NULL,
	created_at 	TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

	PRIMARY KEY (id),
	UNIQUE 	KEY (email)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
	`)
	return err
}

func (m *MySQLStorage) createProjectsTable() error {
	_, err := m.db.Exec(`
CREATE TABLE IF NOT EXISTS projects (
	id 					INT UNSIGNED NOT NULL AUTO_INCREMENT,
	name 				VARCHAR(255) NOT NULL,
	created_at	TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

	PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
	`)
	return err
}

func (m *MySQLStorage) createTasksTable() error {
	_, err := m.db.Exec(`
CREATE TABLE IF NOT EXISTS tasks (
	id 						INT UNSIGNED NOT NULL AUTO_INCREMENT,
	name 					VARCHAR(255) NOT NULL,
	status 				ENUM('TODO', 'IN_PROGRESS', 'IN_TESTING', 'DONE') NOT NULL DEFAULT 'TODO',
	project_id 		INT UNSIGNED NOT NULL,
	assigned_to 	INT UNSIGNED NOT NULL,
	created_at 		TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

	PRIMARY KEY (id),
	FOREIGN KEY (project_id) REFERENCES projects(id),
	FOREIGN KEY (assigned_to) REFERENCES users(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
	`)
	return err
}
