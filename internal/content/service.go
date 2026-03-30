package content

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"time"

	"GoHeadless/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service interface {
	CreateRecord(ctx context.Context, collSlug string, record domain.Record) (primitive.ObjectID, error)
	GetRecords(ctx context.Context, collSlug string) ([]domain.Record, error)
	GetRecord(ctx context.Context, collSlug string, id primitive.ObjectID) (domain.Record, error)
	UpdateRecord(ctx context.Context, collSlug string, id primitive.ObjectID, record domain.Record) error
	DeleteRecord(ctx context.Context, collSlug string, id primitive.ObjectID) error
}

type service struct {
	recordRepo     domain.RecordRepository
	collectionRepo domain.CollectionRepository
}

func NewService(recordRepo domain.RecordRepository, collectionRepo domain.CollectionRepository) Service {
	return &service{
		recordRepo:     recordRepo,
		collectionRepo: collectionRepo,
	}
}

func (s *service) validateSchema(coll *domain.Collection, record domain.Record) error {
	for _, field := range coll.Fields {
		val, ok := record[field.Name]
		if !ok {
			if field.Required {
				return fmt.Errorf("field %s is required", field.Name)
			}
			continue
		}

		// Basic type validation
		switch field.Type {
		case domain.TypeString:
			if _, ok := val.(string); !ok {
				return fmt.Errorf("field %s must be string", field.Name)
			}
		case domain.TypeNumber:
			// JSON usually parses numbers as float64, so we check for that
			if _, ok := val.(int64); !ok {
				if f, ok := val.(float64); ok {
					record[field.Name] = int64(f)
				} else {
					return fmt.Errorf("field %s must be int64", field.Name)
				}
			}
		case domain.TypeFloat:
			if _, ok := val.(float64); !ok {
				return fmt.Errorf("field %s must be float64", field.Name)
			}
		case domain.TypeBool:
			if _, ok := val.(bool); !ok {
				return fmt.Errorf("field %s must be boolean", field.Name)
			}
		case domain.TypeImage:
			if _, ok := val.(string); !ok {
				return fmt.Errorf("field %s must be image path string", field.Name)
			}
		case domain.TypeDateTime:
			// For datetime we can accept string and parse to time.Time
			if strVal, ok := val.(string); ok {
				t, err := time.Parse(time.RFC3339, strVal)
				if err != nil {
					return fmt.Errorf("field %s must be valid RFC3339 datetime string", field.Name)
				}
				record[field.Name] = primitive.NewDateTimeFromTime(t)
			} else if dt, ok := val.(primitive.DateTime); ok {
				record[field.Name] = dt
			} else {
				return fmt.Errorf("field %s has invalid datetime type %v", field.Name, reflect.TypeOf(val))
			}
		}
	}
	return nil
}

func (s *service) CreateRecord(ctx context.Context, collSlug string, record domain.Record) (primitive.ObjectID, error) {
	coll, err := s.collectionRepo.FindBySlug(ctx, collSlug)
	if err != nil {
		return primitive.NilObjectID, errors.New("collection not found")
	}

	if err := s.validateSchema(coll, record); err != nil {
		return primitive.NilObjectID, err
	}

	return s.recordRepo.Create(ctx, collSlug, record)
}

func (s *service) GetRecords(ctx context.Context, collSlug string) ([]domain.Record, error) {
	return s.recordRepo.FindAll(ctx, collSlug)
}

func (s *service) GetRecord(ctx context.Context, collSlug string, id primitive.ObjectID) (domain.Record, error) {
	return s.recordRepo.FindByID(ctx, collSlug, id)
}

func (s *service) UpdateRecord(ctx context.Context, collSlug string, id primitive.ObjectID, record domain.Record) error {
	coll, err := s.collectionRepo.FindBySlug(ctx, collSlug)
	if err != nil {
		return errors.New("collection not found")
	}

	if err := s.validateSchema(coll, record); err != nil {
		return err
	}

	return s.recordRepo.Update(ctx, collSlug, id, record)
}

func (s *service) DeleteRecord(ctx context.Context, collSlug string, id primitive.ObjectID) error {
	return s.recordRepo.Delete(ctx, collSlug, id)
}
