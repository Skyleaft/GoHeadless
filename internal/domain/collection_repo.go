package domain

import (
	"context"
)

type CollectionReader interface {
	FindAll(ctx context.Context) ([]Collection, error)
	FindBySlug(ctx context.Context, slug string) (*Collection, error)
}

type CollectionWriter interface {
	Create(ctx context.Context, coll *Collection) error
	Update(ctx context.Context, coll *Collection) error
	Delete(ctx context.Context, slug string) error
}

type CollectionRepository interface {
	CollectionReader
	CollectionWriter
}
