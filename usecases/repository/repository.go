package repository

import "github.com/norby7/questions-rest-api/entities"

type Repository interface {
	Add(entities.Question) error
	Update(entities.Question) error
	Delete(int64) error
	GetAll(int, int) ([]entities.Question, error)
}
