package models

import (
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Feedback struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    *uuid.UUID         `bson:"user_id," json:"user_id"`
	Text      string             `bson:"text" json:"text"`
	Rating    int                `bson:"rating" json:"rating"`
	Checked   bool               `bson:"checked" json:"checked"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}
