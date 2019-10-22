package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ScheduledAction struct {
	Id         primitive.ObjectID `json:"_id" bson:"_id"`
	EntityType string
	EntityId   string
	ActionType string
	Executed   time.Time
	Result     ActionExecutionResult
}

type ActionExecutionResult struct {
	Code    string
	Message string
	Payload string
}
