package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FieldType string

const (
	// 1. Basic Input Components
	TypeTextInput    FieldType = "text"
	TypeTextArea     FieldType = "textarea"
	TypeNumberInput  FieldType = "number"
	TypeEmailInput   FieldType = "email"
	TypePasswordInput FieldType = "password"
	TypePhoneInput    FieldType = "phone"
	TypeURLInput      FieldType = "url"

	// 2. Selection Components
	TypeSelect      FieldType = "select"
	TypeRadio       FieldType = "radio"
	TypeCheckbox    FieldType = "checkbox"
	TypeMultiSelect FieldType = "multiselect"

	// 3. Date & Time Components
	TypeDatePicker     FieldType = "datepicker"
	TypeTimePicker     FieldType = "timepicker"
	TypeDateTimePicker FieldType = "datetimepicker"
	TypeDateRange      FieldType = "daterange"

	// 4. File & Media Components
	TypeFileUpload  FieldType = "file"
	TypeImageUpload FieldType = "image"

	// 5. Advanced Input Components
	TypeBool         FieldType = "bool" // Simple boolean
	TypeToggleField  FieldType = "toggle"
	TypeSliderField  FieldType = "slider"
	TypeRatingField  FieldType = "rating"
	TypeColorPickerField FieldType = "colorpicker"

	// 6. Layout & Structural Components
	TypeSection    FieldType = "section"
	TypeFieldGroup FieldType = "group"
	TypeTabs       FieldType = "tabs"
	TypeGrid       FieldType = "grid"

	// 8. Repeater / Dynamic List
	TypeRepeater FieldType = "repeater"

	// 9. Relational Components
	TypeRelation     FieldType = "relation"
	TypeAutocomplete FieldType = "autocomplete"

	// 10. Action Components
	TypeSubmitButton FieldType = "submit"
	TypeResetButton  FieldType = "reset"
	TypeCustomAction FieldType = "action"
)

type ValidationRules struct {
	Min       *float64 `bson:"min,omitempty" json:"min,omitempty"`
	Max       *float64 `bson:"max,omitempty" json:"max,omitempty"`
	MinLength *int     `bson:"min_length,omitempty" json:"min_length,omitempty"`
	MaxLength *int     `bson:"max_length,omitempty" json:"max_length,omitempty"`
	Regex     string   `bson:"regex,omitempty" json:"regex,omitempty"`
}

type ArrayConfig struct {
	MinItems    *int `bson:"min_items,omitempty" json:"min_items,omitempty"`
	MaxItems    *int `bson:"max_items,omitempty" json:"max_items,omitempty"`
	UniqueItems bool `bson:"unique_items" json:"unique_items"`
}

type ConditionalLogic struct {
	ShowIf *Condition `bson:"show_if,omitempty" json:"show_if,omitempty"`
}

type Condition struct {
	Field  string      `bson:"field" json:"field"`
	Equals interface{} `bson:"equals" json:"equals"`
}

type Option struct {
	Label string      `bson:"label" json:"label"`
	Value interface{} `bson:"value" json:"value"`
}

type RelationConfig struct {
	Collection string `bson:"collection" json:"collection"`
	Field      string `bson:"field" json:"field"` // UI display field
}

type Field struct {
	Key          string            `bson:"key" json:"key" validate:"required"`
	Label        string            `bson:"label" json:"label"`
	Type         FieldType         `bson:"type" json:"type" validate:"required"`
	Placeholder  string            `bson:"placeholder,omitempty" json:"placeholder,omitempty"`
	DefaultValue interface{}       `bson:"default_value,omitempty" json:"default_value,omitempty"`
	Required     bool              `bson:"required" json:"required"`
	Unique       bool              `bson:"unique" json:"unique"`
	Description  string            `bson:"description,omitempty" json:"description,omitempty"`
	Options      []Option          `bson:"options,omitempty" json:"options,omitempty"`
	Validation   *ValidationRules  `bson:"validation,omitempty" json:"validation,omitempty"`
	Logic        *ConditionalLogic `bson:"logic,omitempty" json:"logic,omitempty"`
	Fields       []Field           `bson:"fields,omitempty" json:"fields,omitempty"`

	// Relational & Advanced Logic
	Relation   *RelationConfig `bson:"relation,omitempty" json:"relation,omitempty"`
	ComputedBy string          `bson:"computed_by,omitempty" json:"computed_by,omitempty"`

	// Array Configuration
	IsArray     bool         `bson:"is_array" json:"is_array"`
	ArrayConfig *ArrayConfig `bson:"array_config,omitempty" json:"array_config,omitempty"`

	// Additional dynamic properties for specific components
	Props map[string]interface{} `bson:"props,omitempty" json:"props,omitempty"`

	// Query engine: include field in full-text-style search across records
	Searchable bool `bson:"searchable,omitempty" json:"searchable,omitempty"`
	// Strip from API responses for anonymous access to public collections
	Internal bool `bson:"internal,omitempty" json:"internal,omitempty"`
}

type CRUDPolicy struct {
	Create []string `bson:"create" json:"create"`
	Read   []string `bson:"read" json:"read"`
	Update []string `bson:"update" json:"update"`
	Delete []string `bson:"delete" json:"delete"`
}

type AccessControl struct {
	IsPublic     bool       `bson:"is_public" json:"is_public"`
	AllowedRoles []string   `bson:"allowed_roles,omitempty" json:"allowed_roles,omitempty"`
	CRUDPolicy   CRUDPolicy `bson:"crud_policy,omitempty" json:"crud_policy,omitempty"`
}

type Collection struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name        string             `bson:"name" json:"name" validate:"required"`
	Slug        string             `bson:"slug" json:"slug" validate:"required"`
	Fields      []Field            `bson:"fields" json:"fields" validate:"required,dive"`
	Description string             `bson:"description" json:"description"`
	Access      *AccessControl     `bson:"access,omitempty" json:"access,omitempty"`
}
