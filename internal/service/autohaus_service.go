// Package service enthält die Geschäftslogik für alle Entitäten.
package service

import (
	"context"
	"fmt"

	"github.com/myorg/myservice/internal/model"
	"github.com/myorg/myservice/internal/repository"
)

// AutohausService defines the business-logic operations for Autohaus.
type AutohausService interface {
	GetByID(ctx context.Context, id uint) (*model.Autohaus, error)
	Create(ctx context.Context, a *model.Autohaus) (uint, error)
	Update(ctx context.Context, a *model.Autohaus) error
	Delete(ctx context.Context, id uint) error
}

type autohausService struct {
	repo repository.AutohausRepository
}

// New creates a new AutohausService with the given repository.
func New(repo repository.AutohausRepository) AutohausService {
	return &autohausService{repo: repo}
}

func (s *autohausService) GetByID(ctx context.Context, id uint) (*model.Autohaus, error) {
	ah, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("service.GetByID id=%d: %w", id, err)
	}
	return ah, nil
}

func (s *autohausService) Create(ctx context.Context, a *model.Autohaus) (uint, error) {
	id, err := s.repo.Create(ctx, a)
	if err != nil {
		return 0, fmt.Errorf("service.Create: %w", err)
	}
	return id, nil
}

func (s *autohausService) Update(ctx context.Context, a *model.Autohaus) error {
	if err := s.repo.Update(ctx, a); err != nil {
		return fmt.Errorf("service.Update id=%d: %w", a.ID, err)
	}
	return nil
}

func (s *autohausService) Delete(ctx context.Context, id uint) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("service.Delete id=%d: %w", id, err)
	}
	return nil
}
