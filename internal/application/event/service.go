package appEvent

import (
	"context"
	"vigia-verde-go/internal/domain/event"
)

type Repository interface {
	Create(ctx context.Context, ev *event.Event) (string, error)
	FindAll(ctx context.Context, filter event.ListFilterParams) ([]event.ListEventResponse, int, error)
	FindByID(ctx context.Context, id string) (*event.Event, error)
}

type EventService struct {
	repo Repository
}

func NewService(repo Repository) *EventService {
	return &EventService{
		repo: repo,
	}
}

func (s *EventService) ListAll(ctx context.Context, filter event.ListFilterParams) ([]event.ListEventResponse, int, error) {
	if err := filter.Validate(); err != nil {
		return nil, 0, err
	}
	return s.repo.FindAll(ctx, filter)
}

func (s *EventService) GetByID(ctx context.Context, id string) (*event.Event, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *EventService) Create(ctx context.Context, input event.EventDto) (string, error) {
	ev, err := event.New(input)
	if err != nil {
		return "", err
	}

	id, err := s.repo.Create(ctx, ev)
	if err != nil {
		return "", err
	}

	return id, nil
}
