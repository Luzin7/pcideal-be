package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const healthCheckTimeout = 2 * time.Second

type MongoPinger interface {
	Ping(ctx context.Context, rp *readpref.ReadPref) error
}

type HealthController struct {
	mongo            MongoPinger
	uniqueIndexReady bool
}

func NewHealthController(mongo MongoPinger, uniqueIndexReady bool) *HealthController {
	return &HealthController{
		mongo:            mongo,
		uniqueIndexReady: uniqueIndexReady,
	}
}

func (c *HealthController) Handle(ctx *gin.Context) {
	checkCtx, cancel := context.WithTimeout(ctx.Request.Context(), healthCheckTimeout)
	defer cancel()

	if err := c.mongo.Ping(checkCtx, nil); err != nil {
		ctx.JSON(http.StatusServiceUnavailable, healthResponse("unavailable", c.uniqueIndexReady))
		return
	}

	ctx.JSON(http.StatusOK, healthResponse("ok", c.uniqueIndexReady))
}

func healthResponse(status string, uniqueIndexReady bool) gin.H {
	return gin.H{
		"status":             status,
		"unique_index_ready": uniqueIndexReady,
	}
}
