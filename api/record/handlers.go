package record

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/sentinel-official/health-check/context"
	"github.com/sentinel-official/health-check/database"
	"github.com/sentinel-official/health-check/types"
)

func HandlerGetRecords(ctx *context.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		_, err := NewRequestGetRecords(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, types.NewResponseError(1, err.Error()))
			return
		}

		projection := bson.M{
			"_id":                       0,
			"addr":                      1,
			"config_exchange_error":     1,
			"config_exchange_timestamp": 1,
			"info_fetch_error":          1,
			"info_fetch_timestamp":      1,
			"location_fetch_error":      1,
			"location_fetch_timestamp":  1,
			"status":                    1,
		}
		opts := options.Find().
			SetProjection(projection)

		items, err := database.RecordFindAll(ctx, nil, opts)
		if err != nil {
			c.JSON(http.StatusInternalServerError, types.NewResponseError(2, err.Error()))
			return
		}

		c.JSON(http.StatusOK, types.NewResponseResult(items))
	}
}

func HandlerGetRecord(ctx *context.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		req, err := NewRequestGetRecord(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, types.NewResponseError(1, err.Error()))
			return
		}

		filter := bson.M{
			"addr": req.URI.Addr,
		}
		projection := bson.M{
			"_id":                       0,
			"addr":                      1,
			"config_exchange_error":     1,
			"config_exchange_timestamp": 1,
			"info_fetch_error":          1,
			"info_fetch_timestamp":      1,
			"location_fetch_error":      1,
			"location_fetch_timestamp":  1,
			"status":                    1,
		}
		opts := options.FindOne().
			SetProjection(projection)

		items, err := database.RecordFindOne(ctx, filter, opts)
		if err != nil {
			c.JSON(http.StatusInternalServerError, types.NewResponseError(2, err.Error()))
			return
		}

		c.JSON(http.StatusOK, types.NewResponseResult(items))
	}
}
