package part

import (
	"net/http"
	"time"

	"github.com/Luzin7/pcideal-be/infra/http/middlewares"
	"github.com/Luzin7/pcideal-be/internal/domain/entity"
	"github.com/Luzin7/pcideal-be/internal/errors"
	"github.com/Luzin7/pcideal-be/internal/useCases/buildAttempt"
	"github.com/Luzin7/pcideal-be/internal/useCases/part"
	"github.com/gin-gonic/gin"
)

type GetBuildRecommendationsController struct {
	generateBuildUC     *part.GenerateBuildRecommendationsUseCase
	buildAttemptService *buildAttempt.BuildAttemptService
}

func NewGetBuildRecommendationsController(
	generateBuildUC *part.GenerateBuildRecommendationsUseCase,
	buildAttemptService *buildAttempt.BuildAttemptService,
) *GetBuildRecommendationsController {
	return &GetBuildRecommendationsController{
		generateBuildUC:     generateBuildUC,
		buildAttemptService: buildAttemptService,
	}
}

func (c *GetBuildRecommendationsController) Handle(ctx *gin.Context) {
	clientIP := middlewares.GetClientIP(ctx)

	var req struct {
		UsageType     string `json:"usage_type"`
		CpuPreference string `json:"cpu_preference"`
		GpuPreference string `json:"gpu_preference"`
		Budget        int64  `json:"budget"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Invalid request body"})
		return
	}

	since := time.Now().Add(-1 * time.Hour)
	count, err := c.buildAttemptService.CountBuildAttemptsByIP(clientIP, since)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Internal server error"})
		return
	}

	const maxAttempts = 4
	if count >= maxAttempts {
		ctx.JSON(http.StatusTooManyRequests, gin.H{"code": errors.ErrBuildAttemptLimit().StatusCode, "message": errors.ErrBuildAttemptLimit().Message})
		return
	}

	buildAttempt := &entity.BuildAttempt{
		IP:      clientIP,
		Goal:    req.UsageType,
		Budget:  int64(req.Budget),
		CPUPref: req.CpuPreference,
		GPUPref: req.GpuPreference,
	}
	_ = c.buildAttemptService.CreateBuildAttempt(buildAttempt)

	part, err := c.generateBuildUC.Execute(req.UsageType, req.CpuPreference, req.GpuPreference, req.Budget)
	if err != nil {
		ctx.JSON(err.StatusCode, gin.H{"code": err.StatusCode, "message": err.Message})
		return
	}

	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, gin.H{"data": part})
}
