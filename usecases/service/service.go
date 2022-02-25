package service

import (
	"questions-rest-api/entities"
	"questions-rest-api/usecases/repository"
)

type Service struct {
	Repo repository.Repository
}

// NewService returns a new Service object address
func NewService(r repository.Repository) *Service {
	return &Service{Repo: r}
}

// Create validates the question object and calls the repository to insert the question
func (s *Service) Create(q entities.Question) error {
	if err := q.Validate(); err != nil {
		return err
	}

	return s.Repo.Add(q)
}

// Update validates the question object and calls the repository to update the question
func (s *Service) Update(q entities.Question) error {
	if err := q.Validate(); err != nil {
		return err
	}

	return s.Repo.Update(q)
}

// Remove calls the repository to delete the question with the given id
func (s *Service) Remove(id int64) error {
	return s.Repo.Delete(id)
}

// ListAll calls the repository to return all questions from the database
func (s *Service) ListAll(lastId, size int) ([]entities.Question, error) {
	return s.Repo.GetAll(lastId, size)
}
