package treeeventservice

import (
	"context"

	treeevent "vigia-verde-go/internal/core"
)

type Repository interface {
	Create(ctx context.Context, input treeevent.TreeEvent) (string, error)
}

type TreeEventService struct {
	repo Repository
}

func NewTreeEventService(repo Repository) *TreeEventService {
	return &TreeEventService{
		repo: repo,
	}
}

func (s *TreeEventService) CreateTreeEvent(ctx context.Context, input treeevent.CreateInput) (string, error) {
	ev, err := treeevent.New(input)
	if err != nil {
		return "", err
	}

	return s.repo.Create(ctx, ev)
}
