package models

import "time"

type Notification struct {
	Type      string    `json:"type" bson:"type"`
	UserID    string    `json:"userId" bson:"userId"`
	Message   string    `json:"message" bson:"message"`
	Data      any       `json:"data,omitempty" bson:"data,omitempty"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
}
