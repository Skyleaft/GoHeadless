package domain

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RecordReader interface {
	FindAll(ctx context.Context, collectionName string) ([]Record, error)
	FindByID(ctx context.Context, collectionName string, id primitive.ObjectID) (Record, error)
	CountRecords(ctx context.Context, collectionName string) (int64, error)
}

type RecordWriter interface {
	Create(ctx context.Context, collectionName string, record Record) (primitive.ObjectID, error)
	Update(ctx context.Context, collectionName string, id primitive.ObjectID, record Record) error
	Delete(ctx context.Context, collectionName string, id primitive.ObjectID) error
}

type RecordRepository interface {
	RecordReader
	RecordWriter
}
