package collection

import (
	"context"

	"GoHeadless/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoRepo struct {
	coll *mongo.Collection
}

func NewRepository(db *mongo.Database) domain.CollectionRepository {
	return &mongoRepo{
		coll: db.Collection("system_collections"),
	}
}

func (r *mongoRepo) FindAll(ctx context.Context) ([]domain.Collection, error) {
	var result []domain.Collection
	cursor, err := r.coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	if err := cursor.All(ctx, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (r *mongoRepo) FindBySlug(ctx context.Context, slug string) (*domain.Collection, error) {
	var result domain.Collection
	err := r.coll.FindOne(ctx, bson.M{"slug": slug}).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *mongoRepo) Create(ctx context.Context, coll *domain.Collection) error {
	_, err := r.coll.InsertOne(ctx, coll)
	return err
}

func (r *mongoRepo) Update(ctx context.Context, coll *domain.Collection) error {
	_, err := r.coll.ReplaceOne(ctx, bson.M{"slug": coll.Slug}, coll, options.Replace().SetUpsert(true))
	return err
}

func (r *mongoRepo) Delete(ctx context.Context, slug string) error {
	_, err := r.coll.DeleteOne(ctx, bson.M{"slug": slug})
	return err
}

func (r *mongoRepo) CountCollections(ctx context.Context) (int64, error) {
	return r.coll.CountDocuments(ctx, bson.M{})
}
