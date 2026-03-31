package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

// User represents an administrator or user of the system
type User struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Username       string             `bson:"username" json:"username" validate:"required"`
	Email          string             `bson:"email,omitempty" json:"email,omitempty"`
	PasswordHash   string             `bson:"password_hash" json:"-"`
	RoleID         string             `bson:"role_id" json:"role_id" validate:"required"`
	IsInitialAdmin bool               `bson:"is_initial_admin" json:"is_initial_admin"`
	ActiveStatus   bool               `bson:"active_status" json:"active_status"`
}

// Role defines the access control scope for a user
type Role struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name        string             `bson:"name" json:"name" validate:"required"`
	Description string             `bson:"description,omitempty" json:"description,omitempty"`
	// * means all collections, otherwise list of collection slugs
	Permissions []Permission `bson:"permissions" json:"permissions"`
}

type Permission struct {
	CollectionSlug string   `bson:"collection_slug" json:"collection_slug"`
	Actions        []string `bson:"actions" json:"actions"` // "create", "read", "update", "delete"
}
