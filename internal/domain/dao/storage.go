package dao

import "golang_lesson17/internal/domain/model"

type Storage interface {
	Store(event *model.Event) (uint64, error)
	Get(id uint64) (*model.Event, error)
	Delete(id uint64) error
}

func DefaultStorage() (Storage, error) {
	return NewMapStorage(), nil
}
