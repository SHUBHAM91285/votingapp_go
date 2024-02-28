package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID               primitive.ObjectID `bson:"_id"`
	Name             string             `json:"name" validate:"required"`
	Age              int                `json:"age" validate:"required"`
	Email            string             `json:"email"`
	Mobile           string             `json:"mobile" validate:"required"`
	Address          string             `json:"address" validate:"required"`
	AadharCardNumber int                `json:"aadharCardNumber" unique:"true" validate:"required"`
	Password         string             `json:"password" validate:"required"`
	Role             string             `json:"role" enum:"voter,admin" default:"voter" validate:"required"`
	IsVoted          bool               `json:"isVoted" default:"false"`
}
