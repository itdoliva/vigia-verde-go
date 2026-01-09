package treeeventservice

import (
	"context"

	"vigia-verde-go/internal/core"
)

type Repository interface {
	Create(ctx context.Context, input core.TreeEvent) (string, error)
}

type TreeEventService struct {
	repo Repository
}

func NewTreeEventService(repo Repository) *TreeEventService {
	return &TreeEventService{
		repo: repo,
	}
}

func (s *TreeEventService) Create(ctx context.Context, input core.CreateInput) (string, error) {
	ev, err := core.New(input)
	if err != nil {
		return "", err
	}

	return s.repo.Create(ctx, ev)
}
