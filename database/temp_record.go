package database

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/sentinel-official/health-check/context"
	"github.com/sentinel-official/health-check/models"
)

const (
	TempRecordCollectionName = "temp_records"
)

func TempRecordFindAll(ctx *context.Context, filter bson.M, opts ...*options.FindOptions) ([]*models.Record, error) {
	var v []*models.Record
	if err := FindAll(ctx, ctx.Database().Collection(TempRecordCollectionName), filter, &v, opts...); err != nil {
		return nil, findError(err)
	}

	return v, nil
}

func TempRecordFindOneAndUpdate(ctx *context.Context, filter, update bson.M, opts ...*options.FindOneAndUpdateOptions) (*models.Record, error) {
	var v models.Record
	if err := FindOneAndUpdate(ctx, ctx.Database().Collection(TempRecordCollectionName), filter, update, &v, opts...); err != nil {
		return nil, findOneAndUpdateError(err)
	}

	return &v, nil
}

func TempRecordFindOne(ctx *context.Context, filter bson.M, opts ...*options.FindOneOptions) (*models.Record, error) {
	var v models.Record
	if err := FindOne(ctx, ctx.Database().Collection(TempRecordCollectionName), filter, &v, opts...); err != nil {
		return nil, findOneError(err)
	}

	return &v, nil
}

func TempRecordIndexesCreateMany(ctx *context.Context, models []mongo.IndexModel, opts ...*options.CreateIndexesOptions) ([]string, error) {
	return IndexesCreateMany(ctx, ctx.Database().Collection(TempRecordCollectionName), models, opts...)
}

func TempRecordInsertOne(ctx *context.Context, v *models.Record, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	return InsertOne(ctx, ctx.Database().Collection(TempRecordCollectionName), v, opts...)
}

func TempRecordUpdateMany(ctx *context.Context, filter, update bson.M, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return UpdateMany(ctx, ctx.Database().Collection(TempRecordCollectionName), filter, update, opts...)
}

func TempRecordDeleteMany(ctx *context.Context, filter bson.M, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	return DeleteMany(ctx, ctx.Database().Collection(TempRecordCollectionName), filter, opts...)
}

func TempRecordInsertMany(ctx *context.Context, v []*models.Record, opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error) {
	var items []interface{}
	for i := 0; i < len(v); i++ {
		items = append(items, v[i])
	}

	return InsertMany(ctx, ctx.Database().Collection(TempRecordCollectionName), items, opts...)
}
