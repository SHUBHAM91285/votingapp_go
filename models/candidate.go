package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Candidate struct {
	ID        primitive.ObjectID `bson:"_id"`
	Name      string             `json:"name" validate:"required"`
	Party     string             `json:"party" validate:"required"`
	Age       int                `json:"age" validate:"required"`
	Votes     []Vote             `json:"votes"`
	VoteCount int                `json:"voteCount" default:"0"`
}

type Vote struct {
	UserID  primitive.ObjectID `json:"user" validate:"required"`
	VotedAt time.Time          `json:"votedAt" default:"Date.now()"`
}
