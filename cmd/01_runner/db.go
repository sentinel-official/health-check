package main

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/sentinel-official/health-check/context"
	"github.com/sentinel-official/health-check/database"
)

func createIndexes(ctx *context.Context) error {
	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{
				bson.E{Key: "addr", Value: 1},
			},
			Options: options.Index().
				SetUnique(true),
		},
		{
			Keys: bson.D{
				bson.E{Key: "ip_addr", Value: 1},
			},
		},
	}

	_, err := database.RecordIndexesCreateMany(ctx, indexes)
	if err != nil {
		return err
	}
	_, err = database.TempRecordIndexesCreateMany(ctx, indexes)
	if err != nil {
		return err
	}

	return nil
}
