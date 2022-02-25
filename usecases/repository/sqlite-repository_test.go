package repository

import (
	"database/sql"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/norby7/questions-rest-api/entities"
	"testing"
)

var (
	dbMock sqlmock.Sqlmock
)

var MockOpener = func(string, string) (*sql.DB, error) {
	db, mock, err := sqlmock.New()

	if err != nil {
		return nil, fmt.Errorf("unable to create mock driver: %s", err.Error())
	}

	dbMock = mock

	return db, nil
}

var MockErrOpener = func(d string, p string) (*sql.DB, error) {
	if p == "errPath" {
		return nil, fmt.Errorf("unable to connect to database")
	}

	return nil, nil
}

func TestNewRepository(t *testing.T) {
	SqlOpen = MockErrOpener
	testCases := []struct {
		name    string
		input   string
		isError bool
	}{{
		name:    "valid path",
		input:   "./database/questions.db",
		isError: false,
	}, {
		name:    "invalid path",
		input:   "errPath",
		isError: true,
	}}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewSqliteRepository(tc.input)

			if (err != nil) != tc.isError {
				t.Errorf("expected error (%v), got error (%v)", tc.isError, err.Error())
			}
		})
	}
}

func TestValidAdd(t *testing.T) {
	SqlOpen = MockOpener
	repo, err := NewSqliteRepository("./test.db")
	if err != nil {
		t.Fatalf("unable to create mock repository: %s", err.Error())
	}

	q := entities.Question{
		Body: "Where does the sun set?",
		Options: []entities.Option{{
			QuestionId:  1,
			Body:        "East",
			Correct:     false,
			OptionOrder: 0,
		}, {
			QuestionId:  1,
			Body:        "West",
			Correct:     true,
			OptionOrder: 1,
		}},
	}

	dbMock.ExpectBegin()
	dbMock.ExpectExec(`INSERT INTO questions`).WithArgs(q.Body).WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectCommit()

	var o entities.Option
	dbMock.ExpectBegin()
	o = q.Options[0]
	dbMock.ExpectExec(`INSERT INTO options`).WithArgs(o.QuestionId, o.Body, o.Correct, o.OptionOrder).WillReturnResult(sqlmock.NewResult(1, 1))
	o = q.Options[1]
	dbMock.ExpectExec(`INSERT INTO options`).WithArgs(1, o.Body, o.Correct, o.OptionOrder).WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectCommit()

	err = repo.Add(q)
	if err != nil {
		t.Fatalf("unable to execute add call: %s", err.Error())
	}
}

func TestBeginErrorAdd(t *testing.T) {
	SqlOpen = MockOpener
	repo, err := NewSqliteRepository("./test.db")
	if err != nil {
		t.Fatalf("unable to create mock repository: %s", err.Error())
	}

	q := entities.Question{
		Body: "Where does the sun set?",
		Options: []entities.Option{{
			QuestionId:  1,
			Body:        "East",
			Correct:     false,
			OptionOrder: 0,
		}, {
			QuestionId:  1,
			Body:        "West",
			Correct:     true,
			OptionOrder: 1,
		}},
	}

	beginErr := fmt.Errorf("error executing begin transaction")

	dbMock.ExpectBegin().WillReturnError(beginErr)

	err = repo.Add(q)
	if err == nil {
		t.Errorf("expected error (%v), got error nil", beginErr)
	}
}

func TestQuestionInsertErrorAdd(t *testing.T) {
	SqlOpen = MockOpener
	repo, err := NewSqliteRepository("./test.db")
	if err != nil {
		t.Fatalf("unable to create mock repository: %s", err.Error())
	}

	q := entities.Question{
		Body: "Where does the sun set?",
		Options: []entities.Option{{
			QuestionId:  1,
			Body:        "East",
			Correct:     false,
			OptionOrder: 0,
		}, {
			QuestionId:  1,
			Body:        "West",
			Correct:     true,
			OptionOrder: 1,
		}},
	}

	execErr := fmt.Errorf("error executing insert question")

	dbMock.ExpectBegin()
	dbMock.ExpectExec(`INSERT INTO questions`).WithArgs(q.Body).WillReturnError(execErr)
	dbMock.ExpectRollback()

	err = repo.Add(q)
	if err == nil {
		t.Errorf("expected error (%v), got error nil", execErr)
	}
}

func TestCommitErrorAdd(t *testing.T) {
	SqlOpen = MockOpener
	repo, err := NewSqliteRepository("./test.db")
	if err != nil {
		t.Fatalf("unable to create mock repository: %s", err.Error())
	}

	q := entities.Question{
		Body: "Where does the sun set?",
		Options: []entities.Option{{
			QuestionId:  1,
			Body:        "East",
			Correct:     false,
			OptionOrder: 0,
		}, {
			QuestionId:  1,
			Body:        "West",
			Correct:     true,
			OptionOrder: 1,
		}},
	}

	commitErr := fmt.Errorf("error commiting transaction")

	dbMock.ExpectBegin()
	dbMock.ExpectExec(`INSERT INTO questions`).WithArgs(q.Body).WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectCommit().WillReturnError(commitErr)
	dbMock.ExpectRollback()

	err = repo.Add(q)
	if err == nil {
		t.Errorf("expected error (%v), got error nil", commitErr)
	}
}

func TestBeginOptionErrorAdd(t *testing.T) {
	SqlOpen = MockOpener
	repo, err := NewSqliteRepository("./test.db")
	if err != nil {
		t.Fatalf("unable to create mock repository: %s", err.Error())
	}

	q := entities.Question{
		Body: "Where does the sun set?",
		Options: []entities.Option{{
			QuestionId:  1,
			Body:        "East",
			Correct:     false,
			OptionOrder: 0,
		}, {
			QuestionId:  1,
			Body:        "West",
			Correct:     true,
			OptionOrder: 1,
		}},
	}

	beginErr := fmt.Errorf("error begining transaction")

	dbMock.ExpectBegin()
	dbMock.ExpectExec(`INSERT INTO questions`).WithArgs(q.Body).WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectCommit()

	dbMock.ExpectBegin().WillReturnError(beginErr)
	dbMock.ExpectRollback()

	dbMock.ExpectBegin()
	dbMock.ExpectExec(`DELETE FROM questions WHERE id = ?`).WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectExec(`DELETE FROM options WHERE questionId = ?`).WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectCommit()

	err = repo.Add(q)
	if err == nil {
		t.Errorf("expected error (%v), got error nil", beginErr)
	}
}

func TestInsertOptionErrorAdd(t *testing.T) {
	SqlOpen = MockOpener
	repo, err := NewSqliteRepository("./test.db")
	if err != nil {
		t.Fatalf("unable to create mock repository: %s", err.Error())
	}

	q := entities.Question{
		Body: "Where does the sun set?",
		Options: []entities.Option{{
			QuestionId:  1,
			Body:        "East",
			Correct:     false,
			OptionOrder: 0,
		}, {
			QuestionId:  1,
			Body:        "West",
			Correct:     true,
			OptionOrder: 1,
		}},
	}

	insertErr := fmt.Errorf("error inserting option")

	dbMock.ExpectBegin()
	dbMock.ExpectExec(`INSERT INTO questions`).WithArgs(q.Body).WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectCommit()

	var o entities.Option
	dbMock.ExpectBegin()
	o = q.Options[0]
	dbMock.ExpectExec(`INSERT INTO options`).WithArgs(o.QuestionId, o.Body, o.Correct, o.OptionOrder).WillReturnResult(sqlmock.NewResult(1, 1))
	o = q.Options[1]
	dbMock.ExpectExec(`INSERT INTO options`).WithArgs(1, o.Body, o.Correct, o.OptionOrder).WillReturnError(insertErr)
	dbMock.ExpectRollback()

	dbMock.ExpectBegin()
	dbMock.ExpectExec(`DELETE FROM questions WHERE id = ?`).WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectExec(`DELETE FROM options WHERE questionId = ?`).WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectCommit()

	err = repo.Add(q)
	if err == nil {
		t.Errorf("expected error (%v), got error nil", insertErr)
	}
}

func TestInsertOptionCommitErrorAdd(t *testing.T) {
	SqlOpen = MockOpener
	repo, err := NewSqliteRepository("./test.db")
	if err != nil {
		t.Fatalf("unable to create mock repository: %s", err.Error())
	}

	q := entities.Question{
		Body: "Where does the sun set?",
		Options: []entities.Option{{
			QuestionId:  1,
			Body:        "East",
			Correct:     false,
			OptionOrder: 0,
		}, {
			QuestionId:  1,
			Body:        "West",
			Correct:     true,
			OptionOrder: 1,
		}},
	}

	commitErr := fmt.Errorf("error commiting transaction")

	dbMock.ExpectBegin()
	dbMock.ExpectExec(`INSERT INTO questions`).WithArgs(q.Body).WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectCommit()

	var o entities.Option
	dbMock.ExpectBegin()
	o = q.Options[0]
	dbMock.ExpectExec(`INSERT INTO options`).WithArgs(o.QuestionId, o.Body, o.Correct, o.OptionOrder).WillReturnResult(sqlmock.NewResult(1, 1))
	o = q.Options[1]
	dbMock.ExpectExec(`INSERT INTO options`).WithArgs(1, o.Body, o.Correct, o.OptionOrder).WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectCommit().WillReturnError(commitErr)
	dbMock.ExpectRollback()

	dbMock.ExpectBegin()
	dbMock.ExpectExec(`DELETE FROM questions WHERE id = ?`).WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectExec(`DELETE FROM options WHERE questionId = ?`).WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectCommit()

	err = repo.Add(q)
	if err == nil {
		t.Errorf("expected error (%v), got error nil", commitErr)
	}
}

func TestValidUpdate(t *testing.T) {
	SqlOpen = MockOpener
	repo, err := NewSqliteRepository("./test.db")
	if err != nil {
		t.Fatalf("unable to create mock repository: %s", err.Error())
	}

	q := entities.Question{
		Id:   1,
		Body: "Where does the sun set?",
		Options: []entities.Option{{
			QuestionId:  1,
			Body:        "East",
			Correct:     false,
			OptionOrder: 0,
		}, {
			QuestionId:  1,
			Body:        "West",
			Correct:     true,
			OptionOrder: 1,
		}},
	}
	var o entities.Option

	dbMock.ExpectBegin()
	dbMock.ExpectExec(`UPDATE questions`).WithArgs(q.Body, q.Id).WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectExec(`DELETE FROM options WHERE questionId = ?`).WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))
	o = q.Options[0]
	dbMock.ExpectExec(`INSERT INTO options`).WithArgs(o.QuestionId, o.Body, o.Correct, o.OptionOrder).WillReturnResult(sqlmock.NewResult(1, 1))
	o = q.Options[1]
	dbMock.ExpectExec(`INSERT INTO options`).WithArgs(1, o.Body, o.Correct, o.OptionOrder).WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectCommit()

	//dbMock.ExpectQuery().WillReturnRows(sqlmock.NewRows())

	err = repo.Update(q)
	if err != nil {
		t.Fatalf("unable to execute add call: %s", err.Error())
	}
}

func TestBeginErrorUpdate(t *testing.T) {
	SqlOpen = MockOpener
	repo, err := NewSqliteRepository("./test.db")
	if err != nil {
		t.Fatalf("unable to create mock repository: %s", err.Error())
	}

	q := entities.Question{
		Id:   1,
		Body: "Where does the sun set?",
		Options: []entities.Option{{
			QuestionId:  1,
			Body:        "East",
			Correct:     false,
			OptionOrder: 0,
		}, {
			QuestionId:  1,
			Body:        "West",
			Correct:     true,
			OptionOrder: 1,
		}},
	}
	beginErr := fmt.Errorf("error begining transaction")

	dbMock.ExpectBegin().WillReturnError(beginErr)

	err = repo.Update(q)
	if err == nil {
		t.Errorf("expected error (%v), got error nil", beginErr)
	}
}

func TestQueryErrorUpdate(t *testing.T) {
	SqlOpen = MockOpener
	repo, err := NewSqliteRepository("./test.db")
	if err != nil {
		t.Fatalf("unable to create mock repository: %s", err.Error())
	}

	q := entities.Question{
		Id:   1,
		Body: "Where does the sun set?",
		Options: []entities.Option{{
			QuestionId:  1,
			Body:        "East",
			Correct:     false,
			OptionOrder: 0,
		}, {
			QuestionId: 1,
			Body:       "West",

			Correct:     true,
			OptionOrder: 1,
		}},
	}
	updateErr := fmt.Errorf("error updating questions")

	dbMock.ExpectBegin()
	dbMock.ExpectExec(`UPDATE questions`).WithArgs(q.Body, q.Id).WillReturnError(updateErr)
	dbMock.ExpectRollback()

	err = repo.Update(q)
	if err == nil {
		t.Errorf("expected error (%v), got error nil", updateErr)
	}
}

func TestDeleteErrorUpdate(t *testing.T) {
	SqlOpen = MockOpener
	repo, err := NewSqliteRepository("./test.db")
	if err != nil {
		t.Fatalf("unable to create mock repository: %s", err.Error())
	}

	q := entities.Question{
		Id:   1,
		Body: "Where does the sun set?",
		Options: []entities.Option{{
			QuestionId:  1,
			Body:        "East",
			Correct:     false,
			OptionOrder: 0,
		}, {
			QuestionId:  1,
			Body:        "West",
			Correct:     true,
			OptionOrder: 1,
		}},
	}

	deleteErr := fmt.Errorf("error deleting options")

	dbMock.ExpectBegin()
	dbMock.ExpectExec(`UPDATE questions`).WithArgs(q.Body, q.Id).WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectExec(`DELETE FROM options WHERE questionId = ?`).WithArgs(1).WillReturnError(deleteErr)
	dbMock.ExpectRollback()

	err = repo.Update(q)
	if err == nil {
		t.Errorf("expected error (%v), got error nil", deleteErr)
	}
}

func TestInsertOptionErrorUpdate(t *testing.T) {
	SqlOpen = MockOpener
	repo, err := NewSqliteRepository("./test.db")
	if err != nil {
		t.Fatalf("unable to create mock repository: %s", err.Error())
	}

	q := entities.Question{
		Id:   1,
		Body: "Where does the sun set?",
		Options: []entities.Option{{
			QuestionId:  1,
			Body:        "East",
			Correct:     false,
			OptionOrder: 0,
		}, {
			QuestionId:  1,
			Body:        "West",
			Correct:     true,
			OptionOrder: 1,
		}},
	}
	var o entities.Option
	insertErr := fmt.Errorf("error inserting options")

	dbMock.ExpectBegin()
	dbMock.ExpectExec(`UPDATE questions`).WithArgs(q.Body, q.Id).WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectExec(`DELETE FROM options WHERE questionId = ?`).WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))
	o = q.Options[0]
	dbMock.ExpectExec(`INSERT INTO options`).WithArgs(o.QuestionId, o.Body, o.Correct, o.OptionOrder).WillReturnResult(sqlmock.NewResult(1, 1))
	o = q.Options[1]
	dbMock.ExpectExec(`INSERT INTO options`).WithArgs(1, o.Body, o.Correct, o.OptionOrder).WillReturnError(insertErr)
	dbMock.ExpectRollback()

	err = repo.Update(q)
	if err == nil {
		t.Errorf("expected error (%v), got error nil", insertErr)
	}
}

func TestCommitErrorUpdate(t *testing.T) {
	SqlOpen = MockOpener
	repo, err := NewSqliteRepository("./test.db")
	if err != nil {
		t.Fatalf("unable to create mock repository: %s", err.Error())
	}

	q := entities.Question{
		Id:   1,
		Body: "Where does the sun set?",
		Options: []entities.Option{{
			QuestionId:  1,
			Body:        "East",
			Correct:     false,
			OptionOrder: 0,
		}, {
			QuestionId:  1,
			Body:        "West",
			Correct:     true,
			OptionOrder: 1,
		}},
	}
	var o entities.Option
	commitErr := fmt.Errorf("error commiting")

	dbMock.ExpectBegin()
	dbMock.ExpectExec(`UPDATE questions`).WithArgs(q.Body, q.Id).WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectExec(`DELETE FROM options WHERE questionId = ?`).WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))
	o = q.Options[0]
	dbMock.ExpectExec(`INSERT INTO options`).WithArgs(o.QuestionId, o.Body, o.Correct, o.OptionOrder).WillReturnResult(sqlmock.NewResult(1, 1))
	o = q.Options[1]
	dbMock.ExpectExec(`INSERT INTO options`).WithArgs(1, o.Body, o.Correct, o.OptionOrder).WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectCommit().WillReturnError(commitErr)
	dbMock.ExpectRollback()

	err = repo.Update(q)
	if err == nil {
		t.Errorf("expected error (%v), got error nil", commitErr)
	}
}

func TestValidGetAll(t *testing.T) {
	SqlOpen = MockOpener
	repo, err := NewSqliteRepository("./test.db")
	if err != nil {
		t.Fatalf("unable to create mock repository: %s", err.Error())
	}

	rows := sqlmock.NewRows([]string{"id", "body"})
	rows.AddRow(1, "Where does the sun set?")
	rows.AddRow(2, "Where does the sun rise?")

	firstOptions := sqlmock.NewRows([]string{"id", "questionId", "body", "correct", "optionOrder"})
	firstOptions.AddRow(1, 1, "West", 0, 0)
	firstOptions.AddRow(2, 1, "East", 1, 0)

	secondOptions := sqlmock.NewRows([]string{"id", "questionId", "body", "correct", "optionOrder"})
	secondOptions.AddRow(3, 2, "West", 1, 0)
	secondOptions.AddRow(4, 2, "East", 0, 0)

	dbMock.ExpectQuery(`SELECT`).WillReturnRows(rows)
	dbMock.ExpectQuery(`SELECT`).WithArgs(1).WillReturnRows(firstOptions)
	dbMock.ExpectQuery(`SELECT`).WithArgs(2).WillReturnRows(secondOptions)

	_, err = repo.GetAll(10, 10)
	if err != nil {
		t.Fatalf("unable to execute get all call: %s", err.Error())
	}
}

func TestQueryQuestionErrorGetAll(t *testing.T) {
	SqlOpen = MockOpener
	repo, err := NewSqliteRepository("./test.db")
	if err != nil {
		t.Fatalf("unable to create mock repository: %s", err.Error())
	}

	queryErr := fmt.Errorf("error fetching data")
	dbMock.ExpectQuery(`SELECT`).WillReturnError(queryErr)

	_, err = repo.GetAll(0, 0)
	if err == nil {
		t.Errorf("expected error (%v), got error nil", queryErr)
	}
}

func TestQueryOptionErrorGetAll(t *testing.T) {
	SqlOpen = MockOpener
	repo, err := NewSqliteRepository("./test.db")
	if err != nil {
		t.Fatalf("unable to create mock repository: %s", err.Error())
	}

	queryErr := fmt.Errorf("error fetching data")

	rows := sqlmock.NewRows([]string{"id", "body"})
	rows.AddRow(1, "Where does the sun set?")
	rows.AddRow(2, "Where does the sun rise?")

	firstOptions := sqlmock.NewRows([]string{"id", "questionId", "body", "correct", "optionOrder"})
	firstOptions.AddRow(1, 1, "West", 0, 0)
	firstOptions.AddRow(2, 1, "East", 1, 0)

	dbMock.ExpectQuery(`SELECT`).WillReturnRows(rows)
	dbMock.ExpectQuery(`SELECT`).WithArgs(1).WillReturnRows(firstOptions)
	dbMock.ExpectQuery(`SELECT`).WithArgs(2).WillReturnError(queryErr)

	_, err = repo.GetAll(0, 0)
	if err == nil {
		t.Errorf("expected error (%v), got error nil", queryErr)
	}
}

func TestValidDelete(t *testing.T) {
	SqlOpen = MockOpener
	repo, err := NewSqliteRepository("./test.db")
	if err != nil {
		t.Fatalf("unable to create mock repository: %s", err.Error())
	}

	dbMock.ExpectBegin()
	dbMock.ExpectExec(`DELETE FROM questions WHERE id = ?`).WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectExec(`DELETE FROM options WHERE questionId = ?`).WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectCommit()

	err = repo.Delete(1)
	if err != nil {
		t.Fatalf("unable to execute delete call: %s", err.Error())
	}
}

func TestCommitErrorDelete(t *testing.T) {
	SqlOpen = MockOpener
	repo, err := NewSqliteRepository("./test.db")
	if err != nil {
		t.Fatalf("unable to create mock repository: %s", err.Error())
	}

	commitErr := fmt.Errorf("error commiting")

	dbMock.ExpectBegin()
	dbMock.ExpectExec(`DELETE FROM questions WHERE id = ?`).WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectExec(`DELETE FROM options WHERE questionId = ?`).WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectCommit().WillReturnError(commitErr)
	dbMock.ExpectRollback()

	err = repo.Delete(1)
	if err != commitErr && err != nil {
		t.Fatalf("unable to execute delete call: %s", err.Error())
	}
}

func TestFirstExecErrorDelete(t *testing.T) {
	SqlOpen = MockOpener
	repo, err := NewSqliteRepository("./test.db")
	if err != nil {
		t.Fatalf("unable to create mock repository: %s", err.Error())
	}

	execErr := fmt.Errorf("error executing delete questions")

	dbMock.ExpectBegin()
	dbMock.ExpectExec(`DELETE FROM questions WHERE id = ?`).WithArgs(1).WillReturnError(execErr)
	dbMock.ExpectRollback()

	err = repo.Delete(1)
	if err != execErr && err != nil {
		t.Fatalf("unable to execute delete call: %s", err.Error())
	}
}

func TestSecondExecErrorDelete(t *testing.T) {
	SqlOpen = MockOpener
	repo, err := NewSqliteRepository("./test.db")
	if err != nil {
		t.Fatalf("unable to create mock repository: %s", err.Error())
	}

	execErr := fmt.Errorf("error executing delete options")

	dbMock.ExpectBegin()
	dbMock.ExpectExec(`DELETE FROM questions WHERE id = ?`).WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectExec(`DELETE FROM options WHERE questionId = ?`).WithArgs(1).WillReturnError(execErr)
	dbMock.ExpectRollback()

	err = repo.Delete(1)
	if err != execErr && err != nil {
		t.Fatalf("unable to execute delete call: %s", err.Error())
	}
}

func TestBeginErrorDelete(t *testing.T) {
	SqlOpen = MockOpener
	repo, err := NewSqliteRepository("./test.db")
	if err != nil {
		t.Fatalf("unable to create mock repository: %s", err.Error())
	}

	beginErr := fmt.Errorf("error executing begin transaction")

	dbMock.ExpectBegin().WillReturnError(beginErr)

	err = repo.Delete(1)
	if err != beginErr && err != nil {
		t.Fatalf("unable to execute delete call: %s", err.Error())
	}
}
