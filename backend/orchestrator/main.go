package main

import (
	"fmt"
	"log"

	"github.com/Davanesh/auto-orchestrator/internal/api"
	"github.com/Davanesh/auto-orchestrator/internal/db"
	"github.com/Davanesh/auto-orchestrator/internal/executors"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	_ = godotenv.Load()

	db.InitDB()

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		AllowCredentials: true,
	}))

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "AutoFlow.AI Orchestrator Running"})
	})

	api.RegisterWorkflowRoutes(r)

	// VENOM webhook integration
r.POST("/webhook/whatsapp", func(c *gin.Context) {
    executors.HandleWhatsAppWebhook(c.Writer, c.Request)
})

	go func() {
		log.Println("🚀 Orchestrator running on port 8080...")
		if err := r.Run(":8080"); err != nil {
			log.Fatal("Server failed:", err)
		}
	}()

	for _, rt := range r.Routes() {
		fmt.Println(rt.Method, rt.Path)
	}

	select {}
}
