package http

import (
	"fmt"
	"github.com/norby7/questions-rest-api/entities"
	"io/ioutil"
	"log"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

type ServiceMock struct {
}

func (s *ServiceMock) Create(u entities.Question) error {
	if u.Body == "errQuestion" {
		return fmt.Errorf("unable to add question")
	}

	return nil
}

func (s *ServiceMock) Update(u entities.Question) error {
	if u.Body == "errQuestion" {
		return fmt.Errorf("unable to update question")
	}

	return nil
}
func (s *ServiceMock) Remove(id int64) error {
	if id != 1 {
		return fmt.Errorf("unable to delete question")
	}

	return nil
}

func (s *ServiceMock) ListAll(lastId, size int) ([]entities.Question, error) {
	if lastId == -2 {
		return []entities.Question{}, fmt.Errorf("error, unable to fetch users")
	}

	return []entities.Question{}, nil
}

func TestAdd(t *testing.T) {
	s := ServiceMock{}
	l := log.New(os.Stdout, "question-api", log.LstdFlags)
	c := NewController(&s, l)

	testCases := []struct {
		name       string
		input      *strings.Reader
		statusCode int
	}{{
		name:       "invalid json object",
		input:      strings.NewReader(`"body":"Where does the sun set?","options":[{"body":"East","correct":false},{"body":"West","correct":true}]}`),
		statusCode: 422,
	}, {
		name:       "add error",
		input:      strings.NewReader(`{"body":"errQuestion","options":[]}`),
		statusCode: 500,
	}, {
		name:       "valid request",
		input:      strings.NewReader(`{"body":"Where does the sun set?","options":[{"body":"East","correct":false},{"body":"West","correct":true}]}`),
		statusCode: 200,
	}}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/question", tc.input)
			rec := httptest.NewRecorder()

			c.Add(rec, req)

			result := rec.Result()

			if tc.statusCode != result.StatusCode {
				resBody, _ := ioutil.ReadAll(result.Body)
				t.Errorf("expected status code (%v), got (%v) with response: (%v)", tc.statusCode, result.StatusCode, string(resBody))
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	s := ServiceMock{}
	l := log.New(os.Stdout, "question-api", log.LstdFlags)
	c := NewController(&s, l)

	testCases := []struct {
		name       string
		id         string
		input      *strings.Reader
		statusCode int
	}{{
		name:       "invalid json object",
		id:         "1",
		input:      strings.NewReader(`"body":"Where does the sun set?","options":[{"body":"East","correct":false},{"body":"West","correct":true}]}`),
		statusCode: 422,
	}, {
		name:       "update error",
		id:         "2",
		input:      strings.NewReader(`{"body":"errQuestion","options":[]}`),
		statusCode: 500,
	}, {
		name:       "missing id",
		id:         "",
		input:      strings.NewReader(`{"body":"errQuestion","options":[]}`),
		statusCode: 400,
	}, {
		name:       "valid request",
		id:         "1",
		input:      strings.NewReader(`{"body":"Where does the sun set?","options":[{"body":"East","correct":false},{"body":"West","correct":true}]}`),
		statusCode: 200,
	}, {
		name:       "zero id",
		id:         "0",
		input:      strings.NewReader(`{"body":"Where does the sun set?","options":[{"body":"East","correct":false},{"body":"West","correct":true}]}`),
		statusCode: 200,
	}}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("PUT", "/question/"+tc.id, tc.input)
			rec := httptest.NewRecorder()

			c.Update(rec, req)
			result := rec.Result()

			if result.StatusCode != tc.statusCode {
				resBody, _ := ioutil.ReadAll(result.Body)
				t.Errorf("expected status code (%v), got (%v) with response: (%v)", tc.statusCode, result.StatusCode, string(resBody))
			}
		})
	}
}

func TestDelete(t *testing.T) {
	s := ServiceMock{}
	l := log.New(os.Stdout, "question-api", log.LstdFlags)
	c := NewController(&s, l)

	testCases := []struct {
		name       string
		input      string
		statusCode int
	}{
		{
			name:       "empty query",
			input:      "",
			statusCode: 400,
		},
		{
			name:       "non integer query",
			input:      "id",
			statusCode: 400,
		},
		{
			name:       "delete error",
			input:      "0",
			statusCode: 500,
		},
		{
			name:       "valid request",
			input:      "1",
			statusCode: 200,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("DELETE", "/question/"+tc.input, nil)
			rec := httptest.NewRecorder()

			c.Delete(rec, req)
			result := rec.Result()

			if result.StatusCode != tc.statusCode {
				resBody, _ := ioutil.ReadAll(result.Body)
				t.Errorf("expected status code (%v), got (%v) with response: (%v)", tc.statusCode, result.StatusCode, string(resBody))
			}
		})
	}
}

func TestGetAll(t *testing.T) {
	s := ServiceMock{}
	l := log.New(os.Stdout, "question-api", log.LstdFlags)
	c := NewController(&s, l)

	testCases := []struct {
		name       string
		input      string
		statusCode int
	}{
		{
			name:       "invalid query parameters structure",
			input:      "?invalidFilter=2question=5",
			statusCode: 200,
		},
		{
			name:       "empty query parameters",
			input:      "",
			statusCode: 200,
		},
		{
			name:       "non empty query parameters",
			input:      "?last_id=10&size=15",
			statusCode: 200,
		},
		{
			name:       "single query parameters",
			input:      "?last_id=10",
			statusCode: 200,
		},
		{
			name:       "negative lastId",
			input:      "?last_id=-1",
			statusCode: 200,
		},
		{
			name:       "get all error",
			input:      "?last_id=-2",
			statusCode: 500,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/questions"+tc.input, nil)
			rec := httptest.NewRecorder()

			c.GetAll(rec, req)
			result := rec.Result()

			if result.StatusCode != tc.statusCode {
				resBody, _ := ioutil.ReadAll(rec.Body)
				t.Errorf("expected status code (%v), got (%v) with response: (%v)", tc.statusCode, result.StatusCode, string(resBody))
			}

		})
	}
}
