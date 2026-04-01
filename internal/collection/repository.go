package collection

import (
	"context"
	"errors"

	"GoHeadless/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoRepo struct {
	db   *mongo.Database
	coll *mongo.Collection
}

func NewRepository(db *mongo.Database) domain.CollectionRepository {
	return &mongoRepo{
		db:   db,
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
	if err != nil {
		return err
	}
	if err := r.ensurePhysicalCollection(ctx, coll.Slug); err != nil {
		_, _ = r.coll.DeleteOne(ctx, bson.M{"slug": coll.Slug})
		return err
	}
	return nil
}

// ensurePhysicalCollection creates the MongoDB collection used by GET /content/{slug}.
// Ignores NamespaceExists (48) so re-runs are safe.
func (r *mongoRepo) ensurePhysicalCollection(ctx context.Context, slug string) error {
	err := r.db.CreateCollection(ctx, slug)
	if err == nil {
		return nil
	}
	var cmdErr mongo.CommandError
	if errors.As(err, &cmdErr) && cmdErr.Code == 48 {
		return nil
	}
	return err
}

func (r *mongoRepo) Update(ctx context.Context, coll *domain.Collection) error {
	_, err := r.coll.ReplaceOne(ctx, bson.M{"slug": coll.Slug}, coll, options.Replace().SetUpsert(true))
	return err
}

func (r *mongoRepo) Delete(ctx context.Context, slug string) error {
	_ = r.db.Collection(slug).Drop(ctx)
	_, err := r.coll.DeleteOne(ctx, bson.M{"slug": slug})
	return err
}

func (r *mongoRepo) CountCollections(ctx context.Context) (int64, error) {
	return r.coll.CountDocuments(ctx, bson.M{})
}

func (r *mongoRepo) EnsurePhysicalCollections(ctx context.Context) error {
	colls, err := r.FindAll(ctx)
	if err != nil {
		return err
	}
	for i := range colls {
		if err := r.ensurePhysicalCollection(ctx, colls[i].Slug); err != nil {
			return err
		}
	}
	return nil
}
