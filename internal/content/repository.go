package content

import (
	"context"

	"GoHeadless/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoRepo struct {
	db *mongo.Database
}

func NewRepository(db *mongo.Database) domain.RecordRepository {
	return &mongoRepo{
		db: db,
	}
}

func (r *mongoRepo) FindAll(ctx context.Context, collectionName string) ([]domain.Record, error) {
	var result []domain.Record
	cursor, err := r.db.Collection(collectionName).Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	if err := cursor.All(ctx, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (r *mongoRepo) FindByID(ctx context.Context, collectionName string, id primitive.ObjectID) (domain.Record, error) {
	var result domain.Record
	err := r.db.Collection(collectionName).FindOne(ctx, bson.M{"_id": id}).Decode(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *mongoRepo) Create(ctx context.Context, collectionName string, record domain.Record) (primitive.ObjectID, error) {
	delete(record, "_id") // Ensure _id is not provided
	res, err := r.db.Collection(collectionName).InsertOne(ctx, record)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return res.InsertedID.(primitive.ObjectID), nil
}

func (r *mongoRepo) Update(ctx context.Context, collectionName string, id primitive.ObjectID, record domain.Record) error {
	delete(record, "_id")
	_, err := r.db.Collection(collectionName).ReplaceOne(ctx, bson.M{"_id": id}, record)
	return err
}

func (r *mongoRepo) Delete(ctx context.Context, collectionName string, id primitive.ObjectID) error {
	_, err := r.db.Collection(collectionName).DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *mongoRepo) CountRecords(ctx context.Context, collectionName string) (int64, error) {
	return r.db.Collection(collectionName).CountDocuments(ctx, bson.M{})
}

func (r *mongoRepo) FindWithOptions(ctx context.Context, collectionName string, filter bson.M, sort bson.D, skip, limit int64, projection bson.M) ([]domain.Record, error) {
	opts := options.Find().SetSort(sort).SetSkip(skip).SetLimit(limit)
	if len(projection) > 0 {
		opts.SetProjection(projection)
	}
	var result []domain.Record
	cursor, err := r.db.Collection(collectionName).Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	if err := cursor.All(ctx, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (r *mongoRepo) CountWithFilter(ctx context.Context, collectionName string, filter bson.M) (int64, error) {
	return r.db.Collection(collectionName).CountDocuments(ctx, filter)
}
