package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/sentinel-official/health-check/models"
)

const (
	SessionCollectionName = "sessions"
)

func SessionFind(ctx context.Context, db *mongo.Database, filter bson.M, opts ...*options.FindOptions) ([]*models.Session, error) {
	var v []*models.Session
	if err := Find(ctx, db.Collection(SessionCollectionName), filter, &v, opts...); err != nil {
		return nil, findError(err)
	}

	return v, nil
}

func SessionFindOneAndUpdate(ctx context.Context, db *mongo.Database, filter, update bson.M, opts ...*options.FindOneAndUpdateOptions) (*models.Session, error) {
	var v models.Session
	if err := FindOneAndUpdate(ctx, db.Collection(SessionCollectionName), filter, update, &v, opts...); err != nil {
		return nil, findOneAndUpdateError(err)
	}

	return &v, nil
}

func SessionFindOne(ctx context.Context, db *mongo.Database, filter bson.M, opts ...*options.FindOneOptions) (*models.Session, error) {
	var v models.Session
	if err := FindOne(ctx, db.Collection(SessionCollectionName), filter, &v, opts...); err != nil {
		return nil, findOneError(err)
	}

	return &v, nil
}

func SessionIndexesCreateMany(ctx context.Context, db *mongo.Database, models []mongo.IndexModel, opts ...*options.CreateIndexesOptions) ([]string, error) {
	return IndexesCreateMany(ctx, db.Collection(SessionCollectionName), models, opts...)
}

func SessionInsertOne(ctx context.Context, db *mongo.Database, v *models.Session, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	return InsertOne(ctx, db.Collection(SessionCollectionName), v, opts...)
}
