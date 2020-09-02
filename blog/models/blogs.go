package models

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (item *BlogItem) Create(coll *mongo.Collection) (*mongo.InsertOneResult, error) {
	return coll.InsertOne(context.Background(), item)
}

func (item *BlogItem) ById(coll *mongo.Collection) (*mongo.SingleResult, error) {
	filter := bson.M{"_id": item.ID}
	res := coll.FindOne(context.Background(), filter)
	err := res.Decode(item)
	if err != nil {
		return nil, err
	}
	return res, err
}

type BlogItem struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	AuthorID string             `bson:"author_id"`
	Content  string             `bson:"content"`
	Title    string             `bson:"title"`
}