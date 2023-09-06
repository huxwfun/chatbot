package storage

import (
	"huxwfun/chatbot/internal/models"
)

type UserStorage struct {
	InMemStorage[models.User]
}

func NewUserStorage() *UserStorage {
	inMem := NewInMemStorage[models.User]()
	return &UserStorage{
		inMem,
	}
}

func (s *UserStorage) GetDefaultBot() {
}
