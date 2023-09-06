package storage

import (
	"context"
)

type InMemStorage[T any] struct {
	storage map[string]T
}

func NewInMemStorage[T any]() InMemStorage[T] {
	storage := map[string]T{}
	return InMemStorage[T]{
		storage: storage,
	}
}

func (s *InMemStorage[T]) Save(ctx context.Context, key string, value T) {
	s.storage[key] = value
}

func (s *InMemStorage[T]) Get(ctx context.Context, key string) (T, bool) {
	value, ok := s.storage[key]
	return value, ok
}

func (s *InMemStorage[T]) GetAll(ctx context.Context) []T {
	values := make([]T, 0, len(s.storage))
	for _, v := range s.storage {
		values = append(values, v)
	}
	return values
}
