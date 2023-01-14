package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DayReview struct {
	ID    primitive.ObjectID `json:"id" bson:"_id"`
	Date  int64              `json:"date"`
	Point int                `json:"point"`
}