package api

import (
	"net/http"
	"github.com/Davanesh/auto-orchestrator/internal/models"
	"github.com/gin-gonic/gin"
)

var workflows []models.Workflow

func GetWorkflows(c *gin.Context) {
	c.JSON(http.StatusOK, workflows)
}

func CreateWorkflow(c *gin.Context) {
	var wf models.Workflow
	if err := c.BindJSON(&wf); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	workflows = append(workflows, wf)
	c.JSON(http.StatusCreated, wf)
}
