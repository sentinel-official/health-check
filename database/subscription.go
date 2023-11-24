package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/sentinel-official/health-check/models"
)

const (
	SubscriptionCollectionName = "subscriptions"
)

func SubscriptionFind(ctx context.Context, db *mongo.Database, filter bson.M, opts ...*options.FindOptions) ([]*models.Subscription, error) {
	var v []*models.Subscription
	if err := Find(ctx, db.Collection(SubscriptionCollectionName), filter, &v, opts...); err != nil {
		return nil, findError(err)
	}

	return v, nil
}

func SubscriptionFindOneAndUpdate(ctx context.Context, db *mongo.Database, filter, update bson.M, opts ...*options.FindOneAndUpdateOptions) (*models.Subscription, error) {
	var v models.Subscription
	if err := FindOneAndUpdate(ctx, db.Collection(SubscriptionCollectionName), filter, update, &v, opts...); err != nil {
		return nil, findOneAndUpdateError(err)
	}

	return &v, nil
}

func SubscriptionFindOne(ctx context.Context, db *mongo.Database, filter bson.M, opts ...*options.FindOneOptions) (*models.Subscription, error) {
	var v models.Subscription
	if err := FindOne(ctx, db.Collection(SubscriptionCollectionName), filter, &v, opts...); err != nil {
		return nil, findOneError(err)
	}

	return &v, nil
}

func SubscriptionIndexesCreateMany(ctx context.Context, db *mongo.Database, models []mongo.IndexModel, opts ...*options.CreateIndexesOptions) ([]string, error) {
	return IndexesCreateMany(ctx, db.Collection(SubscriptionCollectionName), models, opts...)
}

func SubscriptionInsertOne(ctx context.Context, db *mongo.Database, v *models.Subscription, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	return InsertOne(ctx, db.Collection(SubscriptionCollectionName), v, opts...)
}
