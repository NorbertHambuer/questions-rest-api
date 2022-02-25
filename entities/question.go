package entities

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator"
	"io"
)

// Question defines the structure for the question object
// swagger: model
type Question struct {
	// the id for this question
	//
	// required: true
	// min: 1
	Id int64 `json:"-"`
	// the actual question content
	//
	// required: true
	// min: 10
	Body string `json:"body" validate:"required,min=10"`
	// list of possible answers
	//
	// required: true
	// min: 2
	Options []Option `json:"options" validate:"required"`
}

var (
	QuestionOptionsLengthError  = fmt.Errorf("question should have at least 2 options")
	QuestionOptionsCorrectError = fmt.Errorf("there isn't a correct option in the list")
)

// Validate checks and validates each field of the question object based on its definition
func (q *Question) Validate() error {
	validate := validator.New()

	if len(q.Options) < 2 {
		return QuestionOptionsLengthError
	}

	// check if there is at least one correct answer
	isCorrectAnswer := false
	for _, v := range q.Options {
		if v.Correct == true {
			isCorrectAnswer = true
			break
		}
	}

	if !isCorrectAnswer {
		return QuestionOptionsCorrectError
	}

	return validate.Struct(q)
}

// ToJSON serializes the contents of the object to JSON
func (q *Question) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(q)
}

// FromJSON deserializes the JSON into the object
func (q *Question) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(q)
}
