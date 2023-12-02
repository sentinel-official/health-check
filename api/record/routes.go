package record

import (
	"github.com/gin-gonic/gin"

	"github.com/sentinel-official/health-check/context"
)

func RegisterRoutes(ctx *context.Context, router gin.IRouter) {
	router.GET("/records", HandlerGetRecords(ctx))
	router.GET("/records/:addr", HandlerGetRecord(ctx))
}
