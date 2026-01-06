package service

import (
	configs "cmd/auth/config"
	"cmd/auth/internal/store"
	"cmd/auth/pkg/jwt"
)

type UserService interface {
	ValidateCredentials(username, password string) (bool, error)
	GenerateToken(username string) (string, error)
	RefreshToken(token string) (string, error)
}

type userService struct {
	userStore  store.UserStore
	jwtManager *jwt.Manager
}

func NewUserService(cfg *configs.Config) UserService {
	jwtManager := jwt.NewManager(cfg.Secret, cfg.TokenDuration)
	userStore := store.NewInMemoryStore()

	return &userService{
		userStore:  userStore,
		jwtManager: jwtManager,
	}
}

func (s *userService) GenerateToken(username string) (string, error) {
	return s.jwtManager.GenerateToken(username)
}

func (s *userService) RefreshToken(token string) (string, error) {
	return s.jwtManager.RefreshToken(token)
}

func (s *userService) ValidateCredentials(username, password string) (bool, error) {
	user, err := s.userStore.Get(username)
	if err != nil {
		return false, err
	}

	if password != user.Password {
		return false, nil
	}

	return true, nil
}
