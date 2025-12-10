package part

import (
	"net/http"

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
	parts, err := c.useCase.Execute()

	if err != nil {
		ctx.JSON(err.StatusCode, gin.H{"code": err.StatusCode, "message": err.Message})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": parts})
}
