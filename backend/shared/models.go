package shared

import (
    "context"
    "log"
    "time"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

// ---------- Mongo Setup ----------
var DB *mongo.Database

func ConnectMongo(uri, dbName string) {
    client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
    if err != nil {
        log.Fatal("Mongo connection failed:", err)
    }
    DB = client.Database(dbName)
    log.Println("✅ Mongo connected:", dbName)
}

// ---------- Models ----------
type Node struct {
    ID        primitive.ObjectID     `bson:"_id,omitempty" json:"id"`
    Type      string                 `bson:"type" json:"type"`
    Label     string                 `bson:"label" json:"label"`
    Metadata  map[string]interface{} `bson:"metadata" json:"metadata"`
    CanvasID  string                 `bson:"canvasId" json:"canvasId"`
    CreatedAt time.Time              `bson:"createdAt" json:"createdAt"`
}

type Connection struct {
    ID        primitive.ObjectID     `bson:"_id,omitempty" json:"id"`
    From      string                 `bson:"from" json:"from"`
    To        string                 `bson:"to" json:"to"`
    Metadata  map[string]interface{} `bson:"metadata" json:"metadata"`
    CanvasID  string                 `bson:"canvasId" json:"canvasId"`
}

type ExecutionLog struct {
    ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    WorkflowID string             `bson:"workflowId" json:"workflowId"`
    TaskName   string             `bson:"taskName" json:"taskName"`
    Status     string             `bson:"status" json:"status"`
    Timestamp  time.Time          `bson:"timestamp" json:"timestamp"`
}
