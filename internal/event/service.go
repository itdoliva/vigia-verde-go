package event

import (
	"context"
)

type Repository interface {
	Create(ctx context.Context, ev *Event) (string, error)
	FindAll(ctx context.Context, filter ListFilterParams) ([]EventResponse, int, error)
	FindByID(ctx context.Context, id string) (*Event, error)
}

type EventService struct {
	repo Repository
}

func NewService(repo Repository) *EventService {
	return &EventService{
		repo: repo,
	}
}

func (s *EventService) ListAll(ctx context.Context, filter ListFilterParams) ([]EventResponse, int, error) {
	if err := filter.Validate(); err != nil {
		return nil, 0, err
	}
	return s.repo.FindAll(ctx, filter)
}

func (s *EventService) GetByID(ctx context.Context, id string) (*Event, error) {
	return s.repo.FindByID(ctx, id)
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
