package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ReviewNote struct {
	BaseModel `json:",inline"`

	ID        primitive.ObjectID `json:"id" bson:"_id"`
	UserID    primitive.ObjectID `json:"user_id"`
	Point     int                `json:"point"`
	CreatedAt int64              `json:"created_at" bson:"created_at"`
}
