package service

import (
	"errors"
	"fmt"
	"github.com/norby7/questions-rest-api/entities"
	"testing"
)

type RepositoryMock struct {
}

var (
	addError    = fmt.Errorf("unable to add the question")
	updateError = fmt.Errorf("unable to update the question")
	deleteError = fmt.Errorf("unable to delete the question")
	getAllError = fmt.Errorf("unable to fetch questions")
)

func (r *RepositoryMock) Add(u entities.Question) error {
	if u.Body != "Where does the sun set?" {
		return addError
	}
	return nil
}

func (r *RepositoryMock) Update(u entities.Question) error {
	if u.Body != "Where does the sun set?" {
		return updateError
	}

	return nil
}
func (r *RepositoryMock) Delete(id int64) error {
	if id != 1 {
		return deleteError
	}

	return nil
}

func (r *RepositoryMock) GetAll(lastId, size int) ([]entities.Question, error) {
	if lastId == -2 {
		return []entities.Question{}, getAllError
	}

	return []entities.Question{}, nil
}

func TestAdd(t *testing.T) {
	r := &RepositoryMock{}
	s := Service{Repo: r}

	testCases := []struct {
		name          string
		input         entities.Question
		expectedError error
	}{
		{
			name: "valid question, add error",
			input: entities.Question{
				Id:   0,
				Body: "add error question",
				Options: []entities.Option{{
					Id:          0,
					QuestionId:  0,
					Body:        "East",
					Correct:     false,
					OptionOrder: 0,
				}, {
					Id:          0,
					QuestionId:  0,
					Body:        "West",
					Correct:     true,
					OptionOrder: 0,
				}},
			},
			expectedError: addError,
		},
		{
			name: "valid question, valid add",
			input: entities.Question{
				Id:   0,
				Body: "Where does the sun set?",
				Options: []entities.Option{{
					Id:          0,
					QuestionId:  0,
					Body:        "East",
					Correct:     false,
					OptionOrder: 0,
				}, {
					Id:          0,
					QuestionId:  0,
					Body:        "West",
					Correct:     true,
					OptionOrder: 0,
				}},
			},
			expectedError: nil,
		},
		{
			name:          "invalid question structure",
			input:         entities.Question{},
			expectedError: entities.QuestionOptionsLengthError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := s.Create(tc.input)

			if !errors.Is(err, tc.expectedError) {
				t.Errorf("expected error (%v), got error (%v)", tc.expectedError, err.Error())
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	r := &RepositoryMock{}
	s := Service{Repo: r}

	testCases := []struct {
		name          string
		input         entities.Question
		expectedError error
	}{
		{
			name: "valid question, update error",
			input: entities.Question{
				Id:   0,
				Body: "add error question",
				Options: []entities.Option{{
					Id:          0,
					QuestionId:  0,
					Body:        "East",
					Correct:     false,
					OptionOrder: 0,
				}, {
					Id:          0,
					QuestionId:  0,
					Body:        "West",
					Correct:     true,
					OptionOrder: 0,
				}},
			},
			expectedError: updateError,
		},
		{
			name: "valid question, valid update",
			input: entities.Question{
				Id:   0,
				Body: "Where does the sun set?",
				Options: []entities.Option{{
					Id:          0,
					QuestionId:  0,
					Body:        "East",
					Correct:     false,
					OptionOrder: 0,
				}, {
					Id:          0,
					QuestionId:  0,
					Body:        "West",
					Correct:     true,
					OptionOrder: 0,
				}},
			},
			expectedError: nil,
		},
		{
			name:          "invalid question structure",
			input:         entities.Question{},
			expectedError: entities.QuestionOptionsLengthError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := s.Update(tc.input)

			if !errors.Is(err, tc.expectedError) {
				t.Errorf("expected error (%v), got error (%v)", tc.expectedError, err.Error())
			}
		})
	}
}

func TestDelete(t *testing.T) {
	r := &RepositoryMock{}
	s := Service{Repo: r}

	testCases := []struct {
		name          string
		input         int64
		expectedError error
	}{
		{
			name:          "valid id, delete error",
			input:         int64(2),
			expectedError: deleteError,
		},
		{
			name:          "valid id, no error",
			input:         int64(1),
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := s.Remove(tc.input)

			if !errors.Is(err, tc.expectedError) {
				t.Errorf("expected error (%v), got error (%v)", tc.expectedError, err.Error())
			}
		})
	}
}

func TestListAll(t *testing.T) {
	r := &RepositoryMock{}
	s := Service{Repo: r}

	testCases := []struct {
		name          string
		input         [2]int
		expectedError error
	}{
		{
			name:          "valid id, get all error",
			input:         [2]int{-2, 10},
			expectedError: getAllError,
		},
		{
			name:          "valid id, no error",
			input:         [2]int{20, 15},
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := s.ListAll(tc.input[0], tc.input[1])

			if !errors.Is(err, tc.expectedError) {
				t.Errorf("expected error (%v), got error (%v)", tc.expectedError, err.Error())
			}
		})
	}
}
