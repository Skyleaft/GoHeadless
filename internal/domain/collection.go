package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FieldType string

const (
	TypeString   FieldType = "String"
	TypeNumber   FieldType = "Number"
	TypeFloat    FieldType = "Float"
	TypeDateTime FieldType = "DateTime"
	TypeBool     FieldType = "Bool"
)

type Field struct {
	Name        string    `bson:"name" json:"name" validate:"required"`
	Type        FieldType `bson:"type" json:"type" validate:"required,oneof=String Number Float DateTime Bool"`
	Required    bool      `bson:"required" json:"required"`
	Unique      bool      `bson:"unique" json:"unique"`
	Description string    `bson:"description" json:"description"`
}

type Collection struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name        string             `bson:"name" json:"name" validate:"required"`
	Slug        string             `bson:"slug" json:"slug" validate:"required"`
	Fields      []Field            `bson:"fields" json:"fields" validate:"required,dive"`
	Description string             `bson:"description" json:"description"`
}
