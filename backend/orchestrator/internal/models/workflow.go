package models

import "time"

type Workflow struct {
	ID          string    `json:"id" bson:"_id,omitempty"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	Status      string    `json:"status"` // draft, running, completed
	Tasks       []Task    `json:"tasks"`
}

type Task struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	Status string `json:"status"`
}
