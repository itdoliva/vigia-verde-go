package service

import (
	"context"

	treeevent "vigia-verde-go/internal/core"
	"vigia-verde-go/internal/repository"
)

type TreeEventService interface {
	CreateTreeEvent(ctx context.Context, input treeevent.CreateTreeEventInput) (string, error)
}

type treeEventService struct {
	repo repository.TreeEventRepository
}

func NewTreeEventService(repo repository.TreeEventRepository) TreeEventService {
	return &treeEventService{
		repo: repo,
	}
}

func (s *treeEventService) CreateTreeEvent(ctx context.Context, input treeevent.CreateTreeEventInput) (string, error) {
	return s.repo.Create(ctx, input)
}
