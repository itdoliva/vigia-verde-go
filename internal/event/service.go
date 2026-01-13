package event

import (
	"context"
)

type Repository interface {
	Create(ctx context.Context, ev *Event) (string, error)
}

type EventService struct {
	repo Repository
}

func NewService(repo Repository) *EventService {
	return &EventService{
		repo: repo,
	}
}

func (s *EventService) Create(ctx context.Context, input CreateInput) (string, error) {
	ev, err := New(input)
	if err != nil {
		return "", err
	}

	id, err := s.repo.Create(ctx, ev)
	if err != nil {
		return "", err
	}

	return id, nil
}
