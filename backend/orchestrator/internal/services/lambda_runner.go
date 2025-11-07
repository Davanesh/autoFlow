package services

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/Davanesh/auto-orchestrator/internal/models"
)

// SimulateLambda simulates the execution of an AWS Lambda function.
func SimulateLambda(taskName string, workflowID string) models.ExecutionLog {
	time.Sleep(1 * time.Second) // simulate processing delay

	success := rand.Float32() > 0.1 // 90% success rate
	status := "completed"
	desc := fmt.Sprintf("Lambda %s executed successfully.", taskName)

	if !success {
		status = "failed"
		desc = fmt.Sprintf("Lambda %s failed during execution.", taskName)
	}

	return models.ExecutionLog{
		WorkflowID:  workflowID,
		TaskName:    taskName,
		Status:      status,
		Timestamp:   time.Now(),
		Description: desc,
		Details: map[string]interface{}{
			"functionType": "AWS Lambda Simulation",
			"duration":     fmt.Sprintf("%dms", 1000+rand.Intn(500)),
		},
	}
}
