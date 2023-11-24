package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/sentinel-official/health-check/models"
)

const (
	NodeCollectionName = "nodes"
)

func NodeFind(ctx context.Context, db *mongo.Database, filter bson.M, opts ...*options.FindOptions) ([]*models.Node, error) {
	var v []*models.Node
	if err := Find(ctx, db.Collection(NodeCollectionName), filter, &v, opts...); err != nil {
		return nil, findError(err)
	}

	return v, nil
}

func NodeFindOneAndUpdate(ctx context.Context, db *mongo.Database, filter, update bson.M, opts ...*options.FindOneAndUpdateOptions) (*models.Node, error) {
	var v models.Node
	if err := FindOneAndUpdate(ctx, db.Collection(NodeCollectionName), filter, update, &v, opts...); err != nil {
		return nil, findOneAndUpdateError(err)
	}

	return &v, nil
}

func NodeFindOne(ctx context.Context, db *mongo.Database, filter bson.M, opts ...*options.FindOneOptions) (*models.Node, error) {
	var v models.Node
	if err := FindOne(ctx, db.Collection(NodeCollectionName), filter, &v, opts...); err != nil {
		return nil, findOneError(err)
	}

	return &v, nil
}

func NodeIndexesCreateMany(ctx context.Context, db *mongo.Database, models []mongo.IndexModel, opts ...*options.CreateIndexesOptions) ([]string, error) {
	return IndexesCreateMany(ctx, db.Collection(NodeCollectionName), models, opts...)
}

func NodeInsertOne(ctx context.Context, db *mongo.Database, v *models.Node, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	return InsertOne(ctx, db.Collection(NodeCollectionName), v, opts...)
}
