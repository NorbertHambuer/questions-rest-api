// Package http, classification of Question REST API
//
// Documentation for Question API
//
// Schemes: http
// BasePath: /question
// Version: 1.0.0
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
// swagger:meta
package http

import (
	"encoding/json"
	"fmt"
	"github.com/norby7/questions-rest-api/entities"
	"github.com/norby7/questions-rest-api/usecases/service"
	"log"
	"net/http"
	"path"
	"strconv"
)

// Data structure representing a single question
// swagger:response questionResponse
type questionResponse struct {
	// A single question object
	// in body:
	Body entities.Question
}

// Generic error message response
// swagger:response errorResponse
type errorResponse struct {
	Message string `json:"message"`
}

// swagger:response noContent
type noContent struct {
}

// Data structure representing a list of users
// swagger:response questionsListResponse
type questionsListResponse struct {
	// in: body
	Body []entities.Question
}

// swagger:parameters Add Update
type questionParam struct {
	// Question object used for Add or Update
	// Note: the ID field is ignored by add operations
	// in: body
	// required: true
	Body entities.Question
}

type Controller struct {
	Service service.Interactor
	Logger  *log.Logger
}

func NewController(s service.Interactor, l *log.Logger) *Controller {
	return &Controller{Service: s, Logger: l}
}

// swagger:route POST /question question Add
// Creates a new question in the database and then returns it in the response
// responses:
// 200: questionResponse
// 422: errorResponse
// 500: errorResponse

// Add creates a new question in the database and returns it
func (c *Controller) Add(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-type", "application/json")
	c.Logger.Println("Handle Add question")

	var q entities.Question
	err := q.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, fmt.Sprintf("unable to parse question object: %s", err.Error()), http.StatusUnprocessableEntity)
		return
	}

	err = c.Service.Create(q)
	if err != nil {
		http.Error(rw, fmt.Sprintf("unable to add question: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	err = q.ToJSON(rw)
	if err != nil {
		http.Error(rw, fmt.Sprintf("unable to encode response: %s", err.Error()), http.StatusInternalServerError)
		return
	}
}

// swagger:route PUT /question/{id} question Update
// Updates an existing question and returns the updated question in the response
// responses:
// 200: noContent
// 422: errorResponse
// 500: errorResponse

// Update updates an existing question and returns the updated question in response
func (c *Controller) Update(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-type", "application/json")
	c.Logger.Println("Handle Update question")

	id, err := strconv.Atoi(path.Base(r.URL.String()))
	if err != nil {
		http.Error(rw, fmt.Sprintf("invalid question id value: %s", err.Error()), http.StatusBadRequest)
		return
	}

	var q entities.Question
	err = q.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, fmt.Sprintf("unable to parse question object: %s", err.Error()), http.StatusUnprocessableEntity)
		return
	}

	q.Id = int64(id)

	err = c.Service.Update(q)
	if err != nil {
		http.Error(rw, fmt.Sprintf("unable to update question: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	err = q.ToJSON(rw)
	if err != nil {
		http.Error(rw, fmt.Sprintf("unable to encode response: %s", err.Error()), http.StatusInternalServerError)
		return
	}
}

// swagger:route DELETE /question/{id} question Delete
// Deletes a question
// responses:
// 200: noContent
// 400: errorResponse
// 500: errorResponse

// Delete removes a question from the database
func (c *Controller) Delete(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-type", "application/json")
	c.Logger.Println("Handle Delete question")

	id, err := strconv.Atoi(path.Base(r.URL.String()))
	if err != nil {
		http.Error(rw, fmt.Sprintf("invalid question id value: %s", err.Error()), http.StatusBadRequest)
		return
	}

	err = c.Service.Remove(int64(id))
	if err != nil {
		http.Error(rw, fmt.Sprintf("unable to delete question: %s", err.Error()), http.StatusInternalServerError)
		return
	}
}

// swagger:route GET /questions questions GetAll
// Returns a list of all questions in the database
// responses:
// 200: questionsListResponse
// 500: errorResponse

// GetAll returns a list of questions
// It can accept two query parameters:
// - last_id: if this parameter is passed, the data will be filtered using seek pagination
// - size: this parameter determines the number of items on each page when using pagination, defaulted to 10
func (c *Controller) GetAll(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-type", "application/json")
	c.Logger.Println("Handle GetAll questions")

	var err error
	lastId, size := 0, 10

	lastIdParam := r.URL.Query().Get("last_id")
	if lastIdParam != "" {
		lastId, err = strconv.Atoi(lastIdParam)
		if err != nil {
			http.Error(rw, fmt.Sprintf("invalid last_id query parameter: %s", err.Error()), http.StatusBadRequest)
			return
		}
	}

	sizeParam := r.URL.Query().Get("size")
	if sizeParam != "" {
		size, err = strconv.Atoi(sizeParam)
		if err != nil {
			http.Error(rw, fmt.Sprintf("invalid size query parameter: %s", err.Error()), http.StatusBadRequest)
			return
		}
	}

	questions, err := c.Service.ListAll(lastId, size)
	if err != nil {
		http.Error(rw, fmt.Sprintf("unable to fetch questions: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(rw).Encode(questions)
	if err != nil {
		http.Error(rw, fmt.Sprintf("unable to encode questions response: %s", err.Error()), http.StatusUnprocessableEntity)
		return
	}
}
