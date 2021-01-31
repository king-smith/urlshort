package urlshort

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

type RedirectoryDatabase struct {
	db *mongo.Database
}

var redirectoryCollectionName = "redirect"

func NewRedirectoryDatabase(db *mongo.Database) *RedirectoryDatabase {
	svc := RedirectoryDatabase{db}

	return &svc
}

func (r *RedirectoryDatabase) InsertMany(ctx context.Context, v []interface{}) error {

	_, err := r.db.Collection(redirectoryCollectionName).InsertMany(ctx, v)
	if err != nil {
		return err
	}

	return nil
}

func (r *RedirectoryDatabase) Find(ctx context.Context, v interface{}, filter bson.M) error {

	// Find documents based on filter
	cursor, err := r.db.Collection(redirectoryCollectionName).Find(ctx, filter)
	if err != nil {
		return err
	}

	// Unmarshal to provided interface
	err = cursor.All(ctx, v)

	return err
}
