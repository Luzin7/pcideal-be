package controllers

import (
	"log"
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
		c.JSON(err.StatusCode, gin.H{"code": err.StatusCode, "message": err.Message})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": parts})
}

func (pc *PartController) GetPartByID(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "ID não pode ser vazio"})
		return
	}

	part, err := pc.PartService.GetPartByID(id)

	if err != nil {
		log.Println("Erro ao buscar peça:", err)
		c.JSON(err.StatusCode, gin.H{"code": err.StatusCode, "message": err.Message})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": part})
}

func (pc *PartController) GetPartByModel(c *gin.Context) {
	model := c.Param("model")
	part, err := pc.PartService.GetPartByModel(model)
	if err != nil {
		c.JSON(err.StatusCode, gin.H{"code": err.StatusCode, "message": err.Message})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": part})
}

func (pc *PartController) GetBuildRecomendations(c *gin.Context) {
	var req struct {
		UsageType     string `json:"usage_type"`
		CpuPreference string `json:"cpu_preference"`
		GpuPreference string `json:"gpu_preference"`
		Budget        int64  `json:"budget"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Invalid request body"})
		return
	}

	part, err := pc.PartService.GenerateBuildRecomendations(req.UsageType, req.CpuPreference, req.GpuPreference, req.Budget)
	if err != nil {
		c.JSON(err.StatusCode, gin.H{"code": err.StatusCode, "message": err.Message})
		return
	}
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, gin.H{"data": part})
}
