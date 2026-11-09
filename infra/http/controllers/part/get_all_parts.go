package part

import (
	"net/http"
	"strconv"

	"github.com/Luzin7/pcideal-be/internal/domain/repository"
	"github.com/Luzin7/pcideal-be/internal/errors"
	"github.com/Luzin7/pcideal-be/internal/useCases/part"
	"github.com/gin-gonic/gin"
)

type GetAllPartsController struct {
	useCase *part.GetAllPartsUseCase
}

func NewGetAllPartsController(useCase *part.GetAllPartsUseCase) *GetAllPartsController {
	return &GetAllPartsController{useCase: useCase}
}

func (c *GetAllPartsController) Handle(ctx *gin.Context) {
	args, err := parsePagination(ctx)
	if err != nil {
		respondWithError(ctx, err)
		return
	}

	parts, err := c.useCase.Execute(ctx.Request.Context(), args)

	if err != nil {
		respondWithError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": parts})
}

func parsePagination(ctx *gin.Context) (repository.FindAllPartsArgs, *errors.ErrService) {
	limit, err := parseLimit(ctx.Query("limit"))
	if err != nil {
		return repository.FindAllPartsArgs{}, err
	}

	offset, err := parseOffset(ctx.Query("offset"))
	if err != nil {
		return repository.FindAllPartsArgs{}, err
	}

	return repository.FindAllPartsArgs{Limit: limit, Offset: offset}, nil
}

func parseLimit(raw string) (int64, *errors.ErrService) {
	if raw == "" {
		return repository.DefaultPartsPageSize, nil
	}

	value, err := strconv.Atoi(raw)
	if err != nil || value <= 0 {
		return 0, errors.ErrInvalidLimit()
	}

	if int64(value) > repository.MaxPartsPageSize {
		return 0, errors.ErrInvalidLimit()
	}

	return int64(value), nil
}

func parseOffset(raw string) (int64, *errors.ErrService) {
	if raw == "" {
		return 0, nil
	}

	value, err := strconv.Atoi(raw)
	if err != nil || value < 0 {
		return 0, errors.ErrInvalidOffset()
	}

	return int64(value), nil
}
