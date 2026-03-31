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
		validateSingleValue := func(fieldType domain.FieldType, fieldVal interface{}) error {
			switch fieldType {
			case domain.TypeTextInput, domain.TypeTextArea, domain.TypeEmailInput, domain.TypePasswordInput, domain.TypePhoneInput, domain.TypeURLInput, domain.TypeColorPickerField:
				if _, ok := fieldVal.(string); !ok {
					return fmt.Errorf("field %s must be a string", field.Key)
				}
				// String length validation
				if field.Validation != nil {
					str := fieldVal.(string)
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
				switch v := fieldVal.(type) {
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
				if _, ok := fieldVal.(bool); !ok {
					return fmt.Errorf("field %s must be a boolean", field.Key)
				}

			case domain.TypeImageUpload, domain.TypeFileUpload:
				if _, ok := fieldVal.(string); !ok {
					return fmt.Errorf("field %s must be a string (file path)", field.Key)
				}

			case domain.TypeDatePicker, domain.TypeTimePicker, domain.TypeDateTimePicker:
				if strVal, ok := fieldVal.(string); ok {
					_, err := time.Parse(time.RFC3339, strVal)
					if err != nil {
						return fmt.Errorf("field %s must be valid RFC3339 datetime string", field.Key)
					}
				} else if _, ok := fieldVal.(primitive.DateTime); !ok {
					return fmt.Errorf("field %s has invalid datetime type %v", field.Key, reflect.TypeOf(fieldVal))
				}

			case domain.TypeSelect, domain.TypeRadio:
				if _, ok := fieldVal.(string); !ok {
					if _, ok := fieldVal.(float64); !ok {
						return fmt.Errorf("field %s must be a string or number", field.Key)
					}
				}
			case domain.TypeRelation, domain.TypeAutocomplete:
				if strVal, ok := fieldVal.(string); ok {
					if _, err := primitive.ObjectIDFromHex(strVal); err != nil {
						return fmt.Errorf("field %s must be a valid ObjectID hex string", field.Key)
					}
				} else if _, ok := fieldVal.(primitive.ObjectID); !ok {
					return fmt.Errorf("field %s must be an ObjectID", field.Key)
				}
			}
			return nil
		}

		// List / Array Logic
		if field.IsArray {
			items, ok := val.([]interface{})
			if !ok {
				return fmt.Errorf("field %s must be an array", field.Key)
			}

			// Validate constraints
			if field.ArrayConfig != nil {
				if field.ArrayConfig.MinItems != nil && len(items) < *field.ArrayConfig.MinItems {
					return fmt.Errorf("field %s must have at least %d items", field.Key, *field.ArrayConfig.MinItems)
				}
				if field.ArrayConfig.MaxItems != nil && len(items) > *field.ArrayConfig.MaxItems {
					return fmt.Errorf("field %s must have at most %d items", field.Key, *field.ArrayConfig.MaxItems)
				}
				if field.ArrayConfig.UniqueItems {
					seen := make(map[interface{}]bool)
					for _, item := range items {
						if seen[item] {
							return fmt.Errorf("field %s must have unique items", field.Key)
						}
						seen[item] = true
					}
				}
			}

			// Validate each item
			for i, item := range items {
				if err := validateSingleValue(field.Type, item); err != nil {
					return fmt.Errorf("error in %s at index %d: %w", field.Key, i, err)
				}
				// Potential post-processing (e.g. ObjectID/Date conversion)
				if field.Type == domain.TypeRelation || field.Type == domain.TypeAutocomplete {
					if strVal, ok := item.(string); ok {
						if id, err := primitive.ObjectIDFromHex(strVal); err == nil {
							items[i] = id
						}
					}
				} else if field.Type == domain.TypeDatePicker || field.Type == domain.TypeTimePicker || field.Type == domain.TypeDateTimePicker {
					if strVal, ok := item.(string); ok {
						if t, err := time.Parse(time.RFC3339, strVal); err == nil {
							items[i] = primitive.NewDateTimeFromTime(t)
						}
					}
				}
			}
			record[field.Key] = items
		} else {
			// Single value validation
			if err := validateSingleValue(field.Type, val); err != nil {
				return err
			}

			// Single value Post-processing
			if field.Type == domain.TypeRelation || field.Type == domain.TypeAutocomplete {
				if strVal, ok := val.(string); ok {
					if id, err := primitive.ObjectIDFromHex(strVal); err == nil {
						record[field.Key] = id
					}
				}
			} else if field.Type == domain.TypeDatePicker || field.Type == domain.TypeTimePicker || field.Type == domain.TypeDateTimePicker {
				if strVal, ok := val.(string); ok {
					if t, err := time.Parse(time.RFC3339, strVal); err == nil {
						record[field.Key] = primitive.NewDateTimeFromTime(t)
					}
				}
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
