package content

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"regexp"
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
	return s.validateFields(coll.Fields, record)
}

func (s *service) validateFields(fields []domain.Field, record domain.Record) error {
	for _, field := range fields {
		val, ok := record[field.Key]
		if !ok {
			if field.Required {
				return fmt.Errorf("field %s (%s) is required", field.Label, field.Key)
			}
			continue
		}

		// Handle Layout & Structural Components (nested logic)
		if field.Type == domain.TypeFieldGroup || field.Type == domain.TypeSection || field.Type == domain.TypeTabs || field.Type == domain.TypeGrid {
			nestedData, ok := val.(map[string]interface{})
			if !ok {
				return fmt.Errorf("field %s must be an object", field.Key)
			}
			if err := s.validateFields(field.Fields, domain.Record(nestedData)); err != nil {
				return fmt.Errorf("error in nested group %s: %w", field.Key, err)
			}
			continue
		}

		// Handle Repeater / Dynamic List
		if field.Type == domain.TypeRepeater {
			items, ok := val.([]interface{})
			if !ok {
				return fmt.Errorf("field %s must be an array", field.Key)
			}
			for i, item := range items {
				itemData, ok := item.(map[string]interface{})
				if !ok {
					return fmt.Errorf("item %d in repeater %s must be an object", i, field.Key)
				}
				if err := s.validateFields(field.Fields, domain.Record(itemData)); err != nil {
					return fmt.Errorf("error in repeater %s at index %d: %w", field.Key, i, err)
				}
			}
			continue
		}

		// Basic type validation
		switch field.Type {
		case domain.TypeTextInput, domain.TypeTextArea, domain.TypeEmailInput, domain.TypePasswordInput, domain.TypePhoneInput, domain.TypeURLInput, domain.TypeColorPickerField:
			if _, ok := val.(string); !ok {
				return fmt.Errorf("field %s must be a string", field.Key)
			}
			// String length validation
			if field.Validation != nil {
				str := val.(string)
				if field.Validation.MinLength != nil && len(str) < *field.Validation.MinLength {
					return fmt.Errorf("field %s must have at least %d characters", field.Key, *field.Validation.MinLength)
				}
				if field.Validation.MaxLength != nil && len(str) > *field.Validation.MaxLength {
					return fmt.Errorf("field %s must have at most %d characters", field.Key, *field.Validation.MaxLength)
				}
				if field.Validation.Regex != "" {
					matched, err := regexp.MatchString(field.Validation.Regex, str)
					if err != nil {
						return fmt.Errorf("invalid regex in field %s: %w", field.Key, err)
					}
					if !matched {
						return fmt.Errorf("field %s does not match required pattern", field.Key)
					}
				}
			}

		case domain.TypeNumberInput, domain.TypeSliderField, domain.TypeRatingField:
			var num float64
			switch v := val.(type) {
			case float64:
				num = v
			case int:
				num = float64(v)
			case int64:
				num = float64(v)
			default:
				return fmt.Errorf("field %s must be a number", field.Key)
			}

			if field.Validation != nil {
				if field.Validation.Min != nil && num < *field.Validation.Min {
					return fmt.Errorf("field %s must be at least %v", field.Key, *field.Validation.Min)
				}
				if field.Validation.Max != nil && num > *field.Validation.Max {
					return fmt.Errorf("field %s must be at most %v", field.Key, *field.Validation.Max)
				}
			}

		case domain.TypeBool, domain.TypeToggleField:
			if _, ok := val.(bool); !ok {
				return fmt.Errorf("field %s must be a boolean", field.Key)
			}

		case domain.TypeImageUpload, domain.TypeFileUpload:
			// Check for multiple files if applicable
			isMultiple := false
			if field.Props != nil {
				if m, ok := field.Props["multiple"].(bool); ok {
					isMultiple = m
				}
			}

			if isMultiple {
				if _, ok := val.([]interface{}); !ok {
					return fmt.Errorf("field %s must be an array of files", field.Key)
				}
			} else {
				if _, ok := val.(string); !ok {
					return fmt.Errorf("field %s must be a string (file path)", field.Key)
				}
			}

		case domain.TypeDatePicker, domain.TypeTimePicker, domain.TypeDateTimePicker:
			if strVal, ok := val.(string); ok {
				t, err := time.Parse(time.RFC3339, strVal)
				if err != nil {
					return fmt.Errorf("field %s must be valid RFC3339 datetime string", field.Key)
				}
				record[field.Key] = primitive.NewDateTimeFromTime(t)
			} else if dt, ok := val.(primitive.DateTime); ok {
				record[field.Key] = dt
			} else {
				return fmt.Errorf("field %s has invalid datetime type %v", field.Key, reflect.TypeOf(val))
			}

		case domain.TypeSelect, domain.TypeRadio:
			// Could validate against options
			if _, ok := val.(string); !ok {
				if _, ok := val.(float64); !ok {
					return fmt.Errorf("field %s must be a string or number", field.Key)
				}
			}
		case domain.TypeMultiSelect, domain.TypeCheckbox:
			// These usually hold arrays of strings
			if _, ok := val.([]interface{}); !ok {
				return fmt.Errorf("field %s must be an array", field.Key)
			}
		case domain.TypeRelation, domain.TypeAutocomplete:
			// Usually stores ObjectID, if hex string provided convert it
			if strVal, ok := val.(string); ok {
				id, err := primitive.ObjectIDFromHex(strVal)
				if err == nil {
					record[field.Key] = id
				} else {
					return fmt.Errorf("field %s must be a valid ObjectID hex string", field.Key)
				}
			} else if _, ok := val.(primitive.ObjectID); !ok {
				return fmt.Errorf("field %s must be an ObjectID", field.Key)
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
