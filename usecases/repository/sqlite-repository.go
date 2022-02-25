package repository

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"os"
	"questions-rest-api/entities"
)

type SqliteRepository struct {
	Handler *sql.DB
}

var (
	SqlOpen = sql.Open // function that connects to a database and returns a connection handler
)

// NewSqliteRepository connects to a sqlite database and returns a repository object that contains the database connection handler
func NewSqliteRepository(p string) (*SqliteRepository, error) {

	db, err := SqlOpen("sqlite3", p)
	if err != nil {
		return nil, fmt.Errorf("unable to open sqlite database: %s", err.Error())
	}

	return &SqliteRepository{Handler: db}, nil
}

// CreateDatabase checks if the database file exists and creates one if it doesn't
func CreateDatabase(p string) error {
	if _, err := os.Stat(p); errors.Is(err, os.ErrNotExist) {
		if _, err = os.Create(p); err != nil {
			return fmt.Errorf("unable to create sqlite database file: %s", err.Error())
		}
	}

	return nil
}

// ValidateSchema checks if the questions and options schemas exist and creates them if they don't exist
func ValidateSchema(db *sql.DB) error {
	dbExists, err := schemaExists(db)
	if err != nil {
		return fmt.Errorf("unable to check if database schema exists: %s", err.Error())
	}

	if !dbExists {
		err = initializeSchema(db)
		if err != nil {
			return fmt.Errorf("unable to create database schema: %s", err.Error())
		}
	}

	return nil
}

// schemaExists checks if the questions and options tables exist
func schemaExists(handler *sql.DB) (bool, error) {
	var n string
	if err := handler.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='questions';`).Scan(&n); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}

		return false, err
	}

	if err := handler.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='options';`).Scan(&n); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

// initializeSchema reads the schema.sql file and executes the queries inside it
func initializeSchema(handler *sql.DB) error {
	c, err := ioutil.ReadFile("./database/schema.sql")
	if err != nil {
		return fmt.Errorf("unable to open databse schema sql script: %s", err.Error())
	}

	sqlScript := string(c)

	// begin transaction
	tx, err := handler.Begin()
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("unable to start transaction: %s", err.Error())
	}

	// execute insert question statement
	_, err = tx.Exec(sqlScript)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("unable to execute create database queries: %s", err.Error())
	}

	// commit transaction
	if err = tx.Commit(); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("unable to commit transation: %s", err.Error())
	}

	return nil
}

// addOptions inserts all options for a question
func (r *SqliteRepository) addOptions(options []entities.Option, questionId int64) error {
	// begin transaction
	tx, err := r.Handler.Begin()
	if err != nil {
		return fmt.Errorf("unable to start transaction: %s", err.Error())
	}

	for i, o := range options {
		o.QuestionId = questionId
		o.OptionOrder = i

		// execute insert question statement
		_, err = tx.Exec(`INSERT INTO options (questionId, body, correct, optionOrder) VALUES (?, ? , ?, ?)`, o.QuestionId, o.Body, o.Correct, o.OptionOrder)
		if err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("unable to execute insert option statement: %s", err.Error())
		}
	}

	// commit transaction
	if err = tx.Commit(); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("unable to commit transation: %s", err.Error())
	}

	return nil
}

// Add inserts a new question into the database and returns an error in case something went wrong
func (r *SqliteRepository) Add(q entities.Question) error {
	// begin transaction
	tx, err := r.Handler.Begin()
	if err != nil {
		return fmt.Errorf("unable to start transaction: %s", err.Error())
	}

	// execute insert question statement
	res, err := tx.Exec(`INSERT INTO questions (body) VALUES (?)`, q.Body)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("unable to execute insert question statement: %s", err.Error())
	}

	// commit transaction
	if err = tx.Commit(); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("unable to commit transation: %s", err.Error())
	}

	// get new question id
	id, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("unable to get last inserted id: %s", err.Error())
	}

	err = r.addOptions(q.Options, id)
	if err != nil {
		// if the options couldn't be inserted, then delete the new question
		if err := r.Delete(id); err != nil {
			return fmt.Errorf("unable to insert question options and to delete question: %s", err.Error())
		}

		return fmt.Errorf("unable to insert question options: %s", err.Error())
	}

	return nil
}

// Update inserts a new question into the database and returns an error in case something went wrong
func (r *SqliteRepository) Update(q entities.Question) error {
	// begin transaction
	tx, err := r.Handler.Begin()
	if err != nil {
		return fmt.Errorf("unable to start transaction: %s", err.Error())
	}

	// execute update question statement
	_, err = tx.Exec(`UPDATE questions SET body = ? WHERE id = ?`, q.Body, q.Id)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("unable to execute update question statement: %s", err.Error())
	}

	// delete old options
	_, err = tx.Exec(`DELETE FROM options WHERE questionId = ?`, q.Id)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	// insert new options
	for i, o := range q.Options {
		o.QuestionId = q.Id
		o.OptionOrder = i

		// execute insert question statement
		_, err = tx.Exec(`INSERT INTO options (questionId, body, correct, optionOrder) VALUES (?, ? , ?, ?)`, o.QuestionId, o.Body, o.Correct, o.OptionOrder)
		if err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("unable to execute insert option statement: %s", err.Error())
		}
	}

	// commit transaction
	if err = tx.Commit(); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("unable to commit transation: %s", err.Error())
	}

	return nil
}

// Delete removes a question from the database and all its options
func (r *SqliteRepository) Delete(id int64) error {
	tx, err := r.Handler.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(`DELETE FROM questions WHERE id = ?`, id)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	_, err = tx.Exec(`DELETE FROM options WHERE questionId = ?`, id)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	if err = tx.Commit(); err != nil {
		_ = tx.Rollback()
		return err
	}

	return nil
}

// getQuestionOptions returns a list of options for the given question ID
func (r *SqliteRepository) getQuestionOptions(id int64) ([]entities.Option, error) {
	rows, err := r.Handler.Query(`SELECT * FROM options WHERE questionId = ? ORDER BY optionOrder`, id)
	if err != nil {
		return nil, fmt.Errorf("unable to query database for question options: %s", err.Error())
	}

	var ol []entities.Option
	for rows.Next() {
		var o entities.Option

		if err = rows.Scan(&o.Id, &o.QuestionId, &o.Body, &o.Correct, &o.OptionOrder); err != nil {
			return nil, fmt.Errorf("unable to scan option row: %s", err.Error())
		}

		ol = append(ol, o)
	}

	return ol, nil
}

// GetAll returns all the questions from the database filtered by the given parameters
func (r *SqliteRepository) GetAll(lastId, size int) ([]entities.Question, error) {
	query := `SELECT * FROM questions`
	if lastId != 0 {
		query = fmt.Sprintf("%s WHERE id < %d ORDER BY id DESC LIMIT %d", query, lastId, size)
	}

	rows, err := r.Handler.Query(query)
	if err != nil {
		return nil, fmt.Errorf("unable to query database: %s", err.Error())
	}

	defer rows.Close()

	ql := []entities.Question{}
	for rows.Next() {
		var q entities.Question

		if err = rows.Scan(&q.Id, &q.Body); err != nil {
			return nil, fmt.Errorf("unable to scan question row: %s", err.Error())
		}

		q.Options, err = r.getQuestionOptions(q.Id)
		if err != nil {
			return nil, err
		}

		ql = append(ql, q)
	}

	return ql, nil
}
