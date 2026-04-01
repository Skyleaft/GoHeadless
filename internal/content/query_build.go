package content

import (
	"fmt"
	"regexp"
	"strings"

	"GoHeadless/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// schemaField describes a leaf field path for filtering and coercion.
type schemaField struct {
	Path       string
	Type       domain.FieldType
	Searchable bool
	Internal   bool
}

func flattenSchemaFields(fields []domain.Field, prefix string) []schemaField {
	var out []schemaField
	for _, f := range fields {
		path := f.Key
		if prefix != "" {
			path = prefix + "." + f.Key
		}

		switch f.Type {
		case domain.TypeFieldGroup, domain.TypeSection, domain.TypeTabs, domain.TypeGrid:
			out = append(out, flattenSchemaFields(f.Fields, path)...)
		case domain.TypeRepeater:
			// Skip repeater subfields for dynamic filters (complex paths)
			continue
		default:
			out = append(out, schemaField{
				Path:       path,
				Type:       f.Type,
				Searchable: f.Searchable,
				Internal:   f.Internal,
			})
		}
	}
	return out
}

func schemaFieldMap(fields []schemaField) map[string]schemaField {
	m := make(map[string]schemaField, len(fields))
	for _, f := range fields {
		m[f.Path] = f
	}
	return m
}

func buildFilterBSON(pq ParsedQuery, schema map[string]schemaField) (bson.M, error) {
	root := bson.M{}

	for _, pf := range pq.Filters {
		sf, ok := schema[pf.Field]
		if !ok {
			return nil, fmt.Errorf("unknown filter field: %s", pf.Field)
		}
		if sf.Internal {
			return nil, fmt.Errorf("cannot filter on internal field: %s", pf.Field)
		}

		coerced, err := coerceFilterValue(sf, pf.Op, pf.Value)
		if err != nil {
			return nil, fmt.Errorf("field %s: %w", pf.Field, err)
		}

		if pf.Op == "eq" || pf.Op == "contains" {
			if _, exists := root[pf.Field]; exists {
				return nil, fmt.Errorf("conflicting filters for field %s", pf.Field)
			}
			root[pf.Field] = coerced
			continue
		}

		existing, _ := root[pf.Field].(bson.M)
		if existing == nil {
			if _, exists := root[pf.Field]; exists {
				return nil, fmt.Errorf("conflicting filters for field %s", pf.Field)
			}
			existing = bson.M{}
		}
		mongoOp, err := mapOpToMongo(pf.Op)
		if err != nil {
			return nil, err
		}
		if _, dup := existing[mongoOp]; dup {
			return nil, fmt.Errorf("duplicate operator for field %s", pf.Field)
		}
		existing[mongoOp] = coerced
		root[pf.Field] = existing
	}

	return root, nil
}

func mapOpToMongo(op string) (string, error) {
	switch op {
	case "eq":
		return "$eq", nil
	case "ne":
		return "$ne", nil
	case "gt":
		return "$gt", nil
	case "gte":
		return "$gte", nil
	case "lt":
		return "$lt", nil
	case "lte":
		return "$lte", nil
	case "in":
		return "$in", nil
	case "nin":
		return "$nin", nil
	default:
		return "", fmt.Errorf("unsupported op: %s", op)
	}
}

func coerceFilterValue(sf schemaField, op string, raw string) (interface{}, error) {
	raw = strings.TrimSpace(raw)
	if op == "in" || op == "nin" {
		parts := splitCSV(raw)
		arr := make([]interface{}, 0, len(parts))
		for _, p := range parts {
			v, err := coerceSingleValue(sf, strings.TrimSpace(p))
			if err != nil {
				return nil, err
			}
			arr = append(arr, v)
		}
		return arr, nil
	}
	if op == "contains" {
		if !isStringLike(sf.Type) {
			return nil, fmt.Errorf("contains only applies to string fields")
		}
		return primitive.Regex{Pattern: regexp.QuoteMeta(raw), Options: "i"}, nil
	}
	return coerceSingleValue(sf, raw)
}

func splitCSV(s string) []string {
	if s == "" {
		return nil
	}
	return strings.Split(s, ",")
}

func isStringLike(ft domain.FieldType) bool {
	switch ft {
	case domain.TypeTextInput, domain.TypeTextArea, domain.TypeEmailInput, domain.TypeURLInput,
		domain.TypePhoneInput, domain.TypePasswordInput, domain.TypeSelect, domain.TypeRadio,
		domain.TypeColorPickerField:
		return true
	default:
		return false
	}
}

func coerceSingleValue(sf schemaField, raw string) (interface{}, error) {
	switch sf.Type {
	case domain.TypeNumberInput, domain.TypeSliderField, domain.TypeRatingField:
		var x float64
		_, err := fmt.Sscanf(raw, "%f", &x)
		if err != nil {
			return nil, fmt.Errorf("expected number")
		}
		return x, nil

	case domain.TypeBool, domain.TypeToggleField:
		switch strings.ToLower(raw) {
		case "true", "1", "yes":
			return true, nil
		case "false", "0", "no":
			return false, nil
		default:
			return nil, fmt.Errorf("expected boolean")
		}

	case domain.TypeDatePicker, domain.TypeTimePicker, domain.TypeDateTimePicker:
		t, err := parseDateTime(raw)
		if err != nil {
			return nil, fmt.Errorf("expected a valid date/datetime string")
		}
		return primitive.NewDateTimeFromTime(t), nil

	case domain.TypeRelation, domain.TypeAutocomplete:
		id, err := primitive.ObjectIDFromHex(raw)
		if err != nil {
			return nil, fmt.Errorf("expected ObjectID hex")
		}
		return id, nil

	default:
		return raw, nil
	}
}

func buildSearchOr(search string, searchable []schemaField) bson.M {
	term := strings.TrimSpace(search)
	if term == "" {
		return bson.M{}
	}
	pattern := primitive.Regex{Pattern: regexp.QuoteMeta(term), Options: "i"}
	var or []bson.M
	for _, sf := range searchable {
		if !sf.Searchable || sf.Internal {
			continue
		}
		switch {
		case isStringLike(sf.Type):
			or = append(or, bson.M{sf.Path: pattern})
		case sf.Type == domain.TypeNumberInput || sf.Type == domain.TypeSliderField || sf.Type == domain.TypeRatingField:
			var x float64
			if _, err := fmt.Sscanf(term, "%f", &x); err == nil {
				or = append(or, bson.M{sf.Path: x})
			}
		case sf.Type == domain.TypeDatePicker || sf.Type == domain.TypeTimePicker || sf.Type == domain.TypeDateTimePicker:
			if t, err := parseDateTime(term); err == nil {
				or = append(or, bson.M{sf.Path: primitive.NewDateTimeFromTime(t)})
			}
		}
	}
	if len(or) == 0 {
		return bson.M{"_id": bson.M{"$exists": false}} // match nothing
	}
	return bson.M{"$or": or}
}

func buildSort(sortField string, sortDesc bool) bson.D {
	if sortField == "" {
		return bson.D{{Key: "_id", Value: -1}}
	}
	dir := 1
	if sortDesc {
		dir = -1
	}
	return bson.D{{Key: sortField, Value: dir}}
}

func mergeFilterAndSearch(baseFilter bson.M, searchDoc bson.M) bson.M {
	if len(baseFilter) == 0 && len(searchDoc) == 0 {
		return bson.M{}
	}
	if len(searchDoc) == 0 {
		return baseFilter
	}
	if len(baseFilter) == 0 {
		return searchDoc
	}
	return bson.M{"$and": bson.A{baseFilter, searchDoc}}
}

func validateSortField(field string, schema map[string]schemaField) error {
	if field == "" {
		return nil
	}
	if field == "_id" {
		return nil
	}
	if _, ok := schema[field]; ok {
		return nil
	}
	return fmt.Errorf("invalid sort field: %s", field)
}

func publicProjection(schema []schemaField) bson.M {
	proj := bson.M{}
	for _, sf := range schema {
		if sf.Internal {
			proj[sf.Path] = 0
		}
	}
	if len(proj) == 0 {
		return nil
	}
	return proj
}
