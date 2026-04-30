package user

import (
	"context"
	"errors"
	"time"
)

var ErrNotFound = errors.New("not found")

type CreateUserInput struct {
	Name   string
	Email  string
	Mobile string
}

type UserService struct {
	repo UserRepository
}

func NewService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Create(ctx context.Context, in CreateUserInput) (*User, error) {
	u := &User{
		Name:      in.Name,
		Email:     in.Email,
		Mobile:    in.Mobile,
		CreatedAt: time.Now().UTC(),
	}
	return s.repo.Create(ctx, u)
}

func (s *UserService) GetByID(ctx context.Context, id int64) (*User, error) {
	return s.repo.GetByID(ctx, id)
}
