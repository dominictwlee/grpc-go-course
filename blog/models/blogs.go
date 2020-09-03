package models

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ListAll(coll *mongo.Collection) ([]BlogItem, error) {
	opts := options.Find().SetSort(bson.D{{"author_id", 1}})
	cursor, err := coll.Find(context.Background(), bson.D{}, opts)
	if err != nil {
		return nil, err
	}

	var results []BlogItem
	if err = cursor.All(context.Background(), &results); err != nil {
		return nil, err
	}
	return results, nil
}

func (item *BlogItem) Update(coll *mongo.Collection) (*BlogItem, error) {
	opts := options.FindOneAndReplace().SetUpsert(false)
	filter := bson.M{"_id": item.ID}
	res := coll.FindOneAndReplace(context.Background(), filter, item, opts)
	if res.Err() != nil {
		return nil, res.Err()
	}

	return item, nil
}

func (item *BlogItem) Create(coll *mongo.Collection) (*mongo.InsertOneResult, error) {
	return coll.InsertOne(context.Background(), item)
}
func (item *BlogItem) Delete(coll *mongo.Collection) (*mongo.DeleteResult, error) {
	filter := bson.D{{"_id", item.ID}}
	res, err := coll.DeleteOne(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	return res, nil

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
