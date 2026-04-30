package order

import (
	"context"
	"errors"
	"time"
)

var ErrNotFound = errors.New("not found")

type CreateOrderInput struct {
	UserID int64
	Total  float64
}

type OrderService struct {
	repo OrderRepository
}

func NewService(repo OrderRepository) *OrderService {
	return &OrderService{repo: repo}
}

func (s *OrderService) Create(ctx context.Context, in CreateOrderInput) (*Order, error) {
	o := &Order{
		UserID:    in.UserID,
		Total:     in.Total,
		CreatedAt: time.Now().UTC(),
	}
	return s.repo.Create(ctx, o)
}

func (s *OrderService) GetByID(ctx context.Context, id int64) (*Order, error) {
	return s.repo.GetByID(ctx, id)
}
