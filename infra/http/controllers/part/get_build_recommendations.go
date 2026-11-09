package part

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/Luzin7/pcideal-be/infra/http/presenters"
	"github.com/Luzin7/pcideal-be/internal/domain/entity"
	"github.com/Luzin7/pcideal-be/internal/dto"
	"github.com/Luzin7/pcideal-be/internal/errors"
	"github.com/gin-gonic/gin"
)

const releaseBuildAttemptTimeout = 2 * time.Second

type BuildRecommendationsGenerator interface {
	Execute(ctx context.Context, args dto.GenerateBuildRecommendationsDTO) (*presenters.RecommendedBuildsPresenter, *errors.ErrService)
}

type BuildAttemptLimiter interface {
	ClaimBuildAttempt(ctx context.Context, buildAttempt *entity.BuildAttempt) *errors.ErrService
	ReleaseBuildAttempt(ctx context.Context, buildAttempt *entity.BuildAttempt) *errors.ErrService
}

type ClientIPResolver interface {
	Resolve(c *gin.Context) string
}

type GetBuildRecommendationsController struct {
	generateBuild BuildRecommendationsGenerator
	buildAttempts BuildAttemptLimiter
	clientIP      ClientIPResolver
}

func NewGetBuildRecommendationsController(
	generateBuild BuildRecommendationsGenerator,
	buildAttempts BuildAttemptLimiter,
	clientIP ClientIPResolver,
) *GetBuildRecommendationsController {
	return &GetBuildRecommendationsController{
		generateBuild: generateBuild,
		buildAttempts: buildAttempts,
		clientIP:      clientIP,
	}
}

func (c *GetBuildRecommendationsController) Handle(ctx *gin.Context) {
	clientIP := c.clientIP.Resolve(ctx)

	var req dto.GenerateBuildRecommendationsDTO
	if err := ctx.ShouldBindJSON(&req); err != nil {
		respondWithError(ctx, errors.ErrInvalidRequestBody())
		return
	}

	if err := req.Validate(); err != nil {
		respondWithError(ctx, err)
		return
	}

	buildAttempt := entity.NewBuildAttempt(
		clientIP, req.UsageType, req.BudgetCents, req.CpuPreference, req.GpuPreference,
	)

	if err := c.buildAttempts.ClaimBuildAttempt(ctx.Request.Context(), buildAttempt); err != nil {
		respondWithError(ctx, err)
		return
	}

	build, err := c.generateBuild.Execute(ctx.Request.Context(), req)
	if err != nil {
		releaseBuildAttempt(ctx, c.buildAttempts, buildAttempt)
		respondWithError(ctx, err)
		return
	}

	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, gin.H{"data": build})
}

func releaseBuildAttempt(ctx *gin.Context, limiter BuildAttemptLimiter, buildAttempt *entity.BuildAttempt) {
	releaseContext, cancelRelease := context.WithTimeout(
		context.WithoutCancel(ctx.Request.Context()),
		releaseBuildAttemptTimeout,
	)
	defer cancelRelease()

	if err := limiter.ReleaseBuildAttempt(releaseContext, buildAttempt); err != nil {
		log.Printf("[GetBuildRecommendations] failed to release build attempt slot: %v", err)
	}
}

func respondWithError(ctx *gin.Context, err *errors.ErrService) {
	ctx.JSON(err.StatusCode, gin.H{"code": err.StatusCode, "message": err.Message})
}
