package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Davanesh/auto-orchestrator/internal/db"
	"github.com/Davanesh/auto-orchestrator/internal/models"
	"github.com/Davanesh/auto-orchestrator/internal/services"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ---------------- ROUTE REGISTRATION ----------------

func RegisterWorkflowRoutes(r *gin.Engine) {
	r.GET("/workflows", GetWorkflows)
	r.POST("/workflows", CreateWorkflow)
	r.PUT("/workflows/:id", UpdateWorkflowStatus)
	r.POST("/workflows/:id/run", RunWorkflow)
	r.PUT("/workflows/:id/structure", SaveWorkflowStructure) // ✅ new route
}

// ---------------- UPDATE STATUS ----------------

func UpdateWorkflowStatus(c *gin.Context) {
	id := c.Param("id")

	var body struct {
		Status string `json:"status"`
	}

	if err := c.BindJSON(&body); err != nil {
		log.Println("❌ Error binding JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("❌ Invalid ID format:", id)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	collection := db.GetCollection("workflows")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": bson.M{"status": body.Status}}

	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Println("❌ Update error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result.MatchedCount == 0 {
		log.Println("⚠️ No workflow found for ID:", id)
		c.JSON(http.StatusNotFound, gin.H{"error": "Workflow not found"})
		return
	}

	log.Println("✅ Workflow status updated:", id, "->", body.Status)

	c.JSON(http.StatusOK, gin.H{
		"message": "Workflow status updated successfully",
		"id":      id,
		"status":  body.Status,
	})
}

// ---------------- GET ALL WORKFLOWS ----------------

func GetWorkflows(c *gin.Context) {
	collection := db.GetCollection("workflows")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println("📦 Fetching all workflows from MongoDB...")

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		log.Println("❌ Error fetching workflows:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var workflows []models.Workflow
	if err := cursor.All(ctx, &workflows); err != nil {
		log.Println("❌ Cursor decode error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("✅ Found %d workflows\n", len(workflows))
	c.JSON(http.StatusOK, workflows)
}

// ---------------- CREATE WORKFLOW ----------------

func CreateWorkflow(c *gin.Context) {
	var wf models.Workflow
	if err := c.BindJSON(&wf); err != nil {
		log.Println("❌ Invalid JSON for workflow:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	wf.CreatedAt = time.Now()
	wf.Status = "draft"

	collection := db.GetCollection("workflows")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println("🧠 Inserting new workflow into MongoDB:", wf.Name)

	result, err := collection.InsertOne(ctx, wf)
	if err != nil {
		log.Println("❌ Insert error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	wf.ID = result.InsertedID.(primitive.ObjectID)

	log.Println("✅ Workflow inserted successfully with ID:", wf.ID.Hex())

	c.JSON(http.StatusCreated, wf)
}

// ---------------- RUN WORKFLOW ----------------

func RunWorkflow(c *gin.Context) {
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workflow ID"})
		return
	}

	collection := db.GetCollection("workflows")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Println("⚙️ Fetching workflow to execute:", id)

	var wf models.Workflow
	if err := collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&wf); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Workflow not found"})
		return
	}

	logCollection := db.GetCollection("execution_logs")
	results := []models.ExecutionLog{}

	for i, node := range wf.Nodes {
		taskName := node.Type
		log.Printf("🚀 [Lambda] Invoking function: %s\n", taskName)

		// Insert "running" log
		runningLog := models.ExecutionLog{
			WorkflowID:  wf.ID.Hex(),
			TaskName:    taskName,
			Status:      "running",
			Timestamp:   time.Now(),
			Description: fmt.Sprintf("Lambda %s started execution.", taskName),
		}
		_, _ = logCollection.InsertOne(ctx, runningLog)

		// Simulate AWS Lambda run
		execLog := services.SimulateLambda(taskName, wf.ID.Hex())
		_, _ = logCollection.InsertOne(ctx, execLog)

		wf.Nodes[i].Status = execLog.Status
		results = append(results, execLog)

		log.Printf("✅ [Lambda] %s: %s\n", execLog.TaskName, execLog.Description)
	}

	// Update overall workflow status
	wf.Status = "completed"
	if _, err := collection.UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{"$set": wf}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Println("🎯 Workflow completed successfully:", wf.ID.Hex())

	c.JSON(http.StatusOK, gin.H{
		"workflowId": wf.ID.Hex(),
		"status":     wf.Status,
		"results":    results,
	})
}


// ---------------- SAVE WORKFLOW STRUCTURE ----------------

func SaveWorkflowStructure(c *gin.Context) {
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workflow ID"})
		return
	}

	var body struct {
		Nodes       []models.Node       `json:"nodes"`
		Connections []models.Connection `json:"connections"`
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	collection := db.GetCollection("workflows")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"nodes":       body.Nodes,
			"connections": body.Connections,
			"updatedAt":   time.Now(),
		},
	}

	_, err = collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Workflow structure updated successfully",
	})
}
