package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ExecutionLog struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	WorkflowID  string            `bson:"workflowId" json:"workflowId"`
	TaskName   string             `bson:"taskName" json:"taskName"`
	Status     string             `bson:"status" json:"status"` // started, running, completed, failed
	Timestamp  time.Time          `bson:"timestamp" json:"timestamp"`
	Description string            `bson:"description,omitempty" json:"description,omitempty"`
	Details    map[string]interface{} `bson:"details,omitempty" json:"details,omitempty"`
}
