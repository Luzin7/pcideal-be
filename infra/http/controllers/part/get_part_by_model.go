package part

import (
	"net/http"

	"github.com/Luzin7/pcideal-be/internal/useCases/part"
	"github.com/gin-gonic/gin"
)

type GetPartByModelController struct {
	useCase *part.GetPartByModelUseCase
}

func NewGetPartByModelController(useCase *part.GetPartByModelUseCase) *GetPartByModelController {
	return &GetPartByModelController{useCase: useCase}
}

func (c *GetPartByModelController) Handle(ctx *gin.Context) {
	model := ctx.Param("model")
	part, err := c.useCase.Execute(model)
	if err != nil {
		ctx.JSON(err.StatusCode, gin.H{"code": err.StatusCode, "message": err.Message})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": part})
}
