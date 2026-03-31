package auth

import (
	"context"
	"errors"

	"GoHeadless/internal/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
	CreateUser(ctx context.Context, user *domain.User) error
	GetUserByID(ctx context.Context, id primitive.ObjectID) (*domain.User, error)
	GetUserByUsername(ctx context.Context, username string) (*domain.User, error)
	CountUsers(ctx context.Context) (int64, error)
	
	CreateRole(ctx context.Context, role *domain.Role) error
	GetRoleByID(ctx context.Context, id primitive.ObjectID) (*domain.Role, error)
	GetRoleByName(ctx context.Context, name string) (*domain.Role, error)
	GetAllRoles(ctx context.Context) ([]domain.Role, error)
	UpdateRole(ctx context.Context, role *domain.Role) error
	DeleteRole(ctx context.Context, id primitive.ObjectID) error
	
	GetAllUsers(ctx context.Context) ([]domain.User, error)
	DeleteUser(ctx context.Context, id primitive.ObjectID) error
	
	CountRoles(ctx context.Context) (int64, error)
}

type repository struct {
	db *mongo.Database
}

func NewRepository(db *mongo.Database) Repository {
	return &repository{db: db}
}

func (r *repository) CreateUser(ctx context.Context, user *domain.User) error {
	col := r.db.Collection("users")
	res, err := col.InsertOne(ctx, user)
	if err != nil {
		return err
	}
	user.ID = res.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *repository) GetUserByID(ctx context.Context, id primitive.ObjectID) (*domain.User, error) {
	col := r.db.Collection("users")
	var user domain.User
	err := col.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil // Not found
		}
		return nil, err
	}
	return &user, nil
}

func (r *repository) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	col := r.db.Collection("users")
	var user domain.User
	err := col.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *repository) CountUsers(ctx context.Context) (int64, error) {
	col := r.db.Collection("users")
	return col.CountDocuments(ctx, bson.M{})
}

func (r *repository) CreateRole(ctx context.Context, role *domain.Role) error {
	col := r.db.Collection("roles")
	res, err := col.InsertOne(ctx, role)
	if err != nil {
		return err
	}
	role.ID = res.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *repository) GetRoleByID(ctx context.Context, id primitive.ObjectID) (*domain.Role, error) {
	col := r.db.Collection("roles")
	var role domain.Role
	err := col.FindOne(ctx, bson.M{"_id": id}).Decode(&role)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &role, nil
}

func (r *repository) GetRoleByName(ctx context.Context, name string) (*domain.Role, error) {
	col := r.db.Collection("roles")
	var role domain.Role
	err := col.FindOne(ctx, bson.M{"name": name}).Decode(&role)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &role, nil
}
func (r *repository) GetAllRoles(ctx context.Context) ([]domain.Role, error) {
	col := r.db.Collection("roles")
	cursor, err := col.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var roles []domain.Role
	if err := cursor.All(ctx, &roles); err != nil {
		return nil, err
	}
	return roles, nil
}

func (r *repository) UpdateRole(ctx context.Context, role *domain.Role) error {
	col := r.db.Collection("roles")
	_, err := col.ReplaceOne(ctx, bson.M{"_id": role.ID}, role)
	return err
}

func (r *repository) DeleteRole(ctx context.Context, id primitive.ObjectID) error {
	col := r.db.Collection("roles")
	_, err := col.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *repository) GetAllUsers(ctx context.Context) ([]domain.User, error) {
	col := r.db.Collection("users")
	cursor, err := col.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []domain.User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}
	return users, nil
}

func (r *repository) DeleteUser(ctx context.Context, id primitive.ObjectID) error {
	col := r.db.Collection("users")
	_, err := col.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *repository) CountRoles(ctx context.Context) (int64, error) {
	col := r.db.Collection("roles")
	return col.CountDocuments(ctx, bson.M{})
}
