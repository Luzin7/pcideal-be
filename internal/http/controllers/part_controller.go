package controllers

import (
	"net/http"

	"github.com/Luzin7/pcideal-be/internal/http/services"
	"github.com/gin-gonic/gin"
)

type PartController struct {
	PartService *services.PartService
}

func NewPartController(partService *services.PartService) *PartController {
	return &PartController{PartService: partService}
}

func (pc *PartController) GetAllParts(c *gin.Context) {
	parts, err := pc.PartService.GetAllParts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar peças"})
		return
	}

	c.JSON(http.StatusOK, parts)
}

func (pc *PartController) GetPartByID(c *gin.Context) {
	id := c.Param("id")
	part, err := pc.PartService.GetPartByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar peça"})
		return
	}

	if part == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Peça não encontrada"})
		return
	}

	c.JSON(http.StatusOK, part)
}
func (pc *PartController) GetPartByModel(c *gin.Context) {
	model := c.Param("model")
	part, err := pc.PartService.GetPartByModel(model)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar peça"})
		return
	}

	if part == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Peça não encontrada"})
		return
	}

	c.JSON(http.StatusOK, part)
}
