package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Review struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	BookID    string             `bson:"book_id" json:"book_id"` // Храним UUID как string
	UserID    string             `bson:"user_id" json:"user_id"` // Храним UUID как string
	Text      string             `bson:"text" json:"text"`
	Rating    int                `bson:"rating" json:"rating"`
	Likes     int                `bson:"likes" json:"likes"`
	Dislikes  int                `bson:"dislikes" json:"dislikes"`
	Versions  []ReviewVersion    `bson:"versions" json:"versions"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

type ReviewVersion struct {
	Text     string    `bson:"text" json:"text"`
	EditedAt time.Time `bson:"edited_at" json:"edited_at"`
	EditedBy string    `bson:"edited_by" json:"edited_by"`
}

type ReviewVote struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ReviewID primitive.ObjectID `bson:"review_id" json:"review_id"`
	UserID   string             `bson:"user_id" json:"user_id"` // Храним UUID как string
	Vote     int                `bson:"vote" json:"vote"`       // 1 - лайк, -1 - дизлайк
}
