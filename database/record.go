package database

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/sentinel-official/health-check/context"
	"github.com/sentinel-official/health-check/models"
)

const (
	RecordCollectionName = "records"
)

func RecordFindAll(ctx *context.Context, filter bson.M, opts ...*options.FindOptions) ([]*models.Record, error) {
	var v []*models.Record
	if err := FindAll(ctx, ctx.Database().Collection(RecordCollectionName), filter, &v, opts...); err != nil {
		return nil, findError(err)
	}

	return v, nil
}

func RecordFindOneAndUpdate(ctx *context.Context, filter, update bson.M, opts ...*options.FindOneAndUpdateOptions) (*models.Record, error) {
	var v models.Record
	if err := FindOneAndUpdate(ctx, ctx.Database().Collection(RecordCollectionName), filter, update, &v, opts...); err != nil {
		return nil, findOneAndUpdateError(err)
	}

	return &v, nil
}

func RecordFindOne(ctx *context.Context, filter bson.M, opts ...*options.FindOneOptions) (*models.Record, error) {
	var v models.Record
	if err := FindOne(ctx, ctx.Database().Collection(RecordCollectionName), filter, &v, opts...); err != nil {
		return nil, findOneError(err)
	}

	return &v, nil
}

func RecordIndexesCreateMany(ctx *context.Context, models []mongo.IndexModel, opts ...*options.CreateIndexesOptions) ([]string, error) {
	return IndexesCreateMany(ctx, ctx.Database().Collection(RecordCollectionName), models, opts...)
}

func RecordInsertOne(ctx *context.Context, v *models.Record, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	return InsertOne(ctx, ctx.Database().Collection(RecordCollectionName), v, opts...)
}

func RecordUpdateMany(ctx *context.Context, filter, update bson.M, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return UpdateMany(ctx, ctx.Database().Collection(RecordCollectionName), filter, update, opts...)
}

func RecordDeleteMany(ctx *context.Context, filter bson.M, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	return DeleteMany(ctx, ctx.Database().Collection(RecordCollectionName), filter, opts...)
}

func RecordInsertMany(ctx *context.Context, v []*models.Record, opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error) {
	var items []interface{}
	for i := 0; i < len(v); i++ {
		items = append(items, v[i])
	}

	return InsertMany(ctx, ctx.Database().Collection(RecordCollectionName), items, opts...)
}
