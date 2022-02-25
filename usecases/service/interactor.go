package service

import "questions-rest-api/entities"

type Interactor interface {
	Create(entities.Question) error
	Update(entities.Question) error
	Remove(int64) error
	ListAll(int, int) ([]entities.Question, error)
}
