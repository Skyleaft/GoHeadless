package domain

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Record map[string]interface{}

func (r Record) ToBSON() (bson.M, error) {
	b, err := bson.Marshal(r)
	if err != nil {
		return nil, err
	}
	var res bson.M
	err = bson.Unmarshal(b, &res)
	return res, err
}

func (r Record) GetID() primitive.ObjectID {
	if val, ok := r["_id"]; ok {
		if id, ok := val.(primitive.ObjectID); ok {
			return id
		}
	}
	return primitive.NilObjectID
}
