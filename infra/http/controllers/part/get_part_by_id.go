package part

import (
	"log"
	"net/http"

	"github.com/Luzin7/pcideal-be/internal/useCases/part"
	"github.com/gin-gonic/gin"
)

type GetPartByIDController struct {
	useCase *part.GetPartByIDUseCase
}

func NewGetPartByIDController(useCase *part.GetPartByIDUseCase) *GetPartByIDController {
	return &GetPartByIDController{useCase: useCase}
}

func (c *GetPartByIDController) Handle(ctx *gin.Context) {
	id := ctx.Param("id")

	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "ID não pode ser vazio"})
		return
	}

	part, err := c.useCase.Execute(id)

	if err != nil {
		log.Println("Erro ao buscar peça:", err)
		ctx.JSON(err.StatusCode, gin.H{"code": err.StatusCode, "message": err.Message})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": part})
}
