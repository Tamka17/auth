package store

import (
	"errors"
)

type UserStore interface {
	Get(username string) (User, error)
}

type User struct {
	Username string
	Password string
}

type InMemoryStore struct {
	Users map[string]User
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		Users: map[string]User{
			"testuser": {
				Username: "testuser",
				Password: "testpass",
			},
		},
	}
}

func (s *InMemoryStore) Get(username string) (User, error) {
	user, exists := s.Users[username]
	if !exists {
		return User{}, errors.New("user not found")
	}
	return user, nil
}
