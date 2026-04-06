package collection

import (
	"context"
	"errors"

	"GoHeadless/internal/domain"
)

type Service interface {
	CreateCollection(ctx context.Context, coll *domain.Collection) error
	GetAllCollections(ctx context.Context) ([]domain.Collection, error)
	GetCollection(ctx context.Context, slug string) (*domain.Collection, error)
	UpdateCollection(ctx context.Context, slug string, coll *domain.Collection) error
	DeleteCollection(ctx context.Context, slug string) error
}

type service struct {
	repo domain.CollectionRepository
}

func NewService(repo domain.CollectionRepository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) CreateCollection(ctx context.Context, coll *domain.Collection) error {
	// Check if already exists
	existing, _ := s.repo.FindBySlug(ctx, coll.Slug)
	if existing != nil {
		return errors.New("collection already exists with this slug")
	}
	return s.repo.Create(ctx, coll)
}

func (s *service) GetAllCollections(ctx context.Context) ([]domain.Collection, error) {
	return s.repo.FindAll(ctx)
}

func (s *service) GetCollection(ctx context.Context, slug string) (*domain.Collection, error) {
	return s.repo.FindBySlug(ctx, slug)
}

func (s *service) UpdateCollection(ctx context.Context, slug string, coll *domain.Collection) error {
	existing, err := s.repo.FindBySlug(ctx, slug)
	if err != nil {
		return errors.New("collection not found")
	}

	// Enforce slug immutability
	if coll.Slug != existing.Slug {
		return errors.New("collection slug cannot be changed")
	}
	
	// Enforce name immutability
	if coll.Name != existing.Name {
		return errors.New("collection name cannot be changed")
	}

	return s.repo.Update(ctx, coll)
}

func (s *service) DeleteCollection(ctx context.Context, slug string) error {
	return s.repo.Delete(ctx, slug)
}
