package controllers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/SHUBHAM91285/votingApp_go/database"
	"github.com/SHUBHAM91285/votingApp_go/models"
	"github.com/SHUBHAM91285/votingApp_go/token"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var candidateCollection *mongo.Collection = database.OpenCollection(database.Client, "candidate")
var validate = validator.New()

func CheckAdminRole(user models.User) bool {
	if user.Role == "admin" {
		return true
	}

	return false
}

func AddCandidate() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var candidate models.Candidate
		var foundUser models.User

		AadharCardNumber := c.Param("AadharCardNumber")
		AadharCard, err := strconv.Atoi(AadharCardNumber)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "internal server error"})
			return
		}
		err = userCollection.FindOne(ctx, bson.M{"aadharcardnumber": AadharCard}).Decode(&foundUser)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
			return
		}
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
		tokenString = tokenString[len("Bearer "):]

		err = token.VerifyToken(tokenString)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}

		adminCheck := CheckAdminRole(foundUser)
		if !adminCheck {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user is not an admin"})
			return
		}

		if err := c.BindJSON(&candidate); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(candidate)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}
		candidate.ID = primitive.NewObjectID()
		candidateInfo, insertErr := candidateCollection.InsertOne(ctx, candidate)
		if insertErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Candidate is not created"})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, gin.H{"userInfo": candidateInfo})

	}
}

func UpdateCandidate() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var candidate models.Candidate
		var foundCandidate models.Candidate
		var foundUser models.User
		AadharCardNumber := c.Param("AadharCardNumber")
		AadharCard, err := strconv.Atoi(AadharCardNumber)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "internal server error"})
			return
		}
		err = userCollection.FindOne(ctx, bson.M{"aadharcardnumber": AadharCard}).Decode(&foundUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
			return
		}
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
		tokenString = tokenString[len("Bearer "):]

		err = token.VerifyToken(tokenString)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}

		adminCheck := CheckAdminRole(foundUser)
		if !adminCheck {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user is not an admin"})
			return
		}
		if err := c.BindJSON(&candidate); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		err = candidateCollection.FindOne(ctx, bson.M{"name": candidate.Name}).Decode(&foundCandidate)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "candidate not found"})
			return
		}

		result, err := candidateCollection.UpdateOne(
			ctx,
			bson.M{"_id": foundCandidate.ID},
			bson.D{{"$set", bson.D{
				{"name", candidate.Name},
				{"party", candidate.Party},
				{"age", candidate.Age},
			}}},
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update candidate"})
			return
		}

		if result.ModifiedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "data not modified"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "candidate updated successfully"})

	}
}

func DeleteCandidate() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var candidate models.Candidate
		var foundUser models.User
		AadharCardNumber := c.Param("AadharCardNumber")
		AadharCard, err := strconv.Atoi(AadharCardNumber)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "internal server error"})
			return
		}
		err = userCollection.FindOne(ctx, bson.M{"aadharcardnumber": AadharCard}).Decode(&foundUser)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
			return
		}
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
		tokenString = tokenString[len("Bearer "):]

		err = token.VerifyToken(tokenString)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}

		adminCheck := CheckAdminRole(foundUser)
		if !adminCheck {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user is not an admin"})
			return
		}
		if err := c.BindJSON(&candidate); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		result, err := candidateCollection.DeleteOne(ctx, bson.M{"name": candidate.Name})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete candidate"})
			return
		}

		if result.DeletedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "candidate not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "candidate deleted successfully"})
	}
}

func VoteCandidate() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var candidate models.Candidate
		var foundUser models.User
		var foundCandidate models.Candidate
		var voteDetails models.Vote

		AadharCardNumber := c.Param("AadharCardNumber")
		AadharCard, err := strconv.Atoi(AadharCardNumber)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "internal server error"})
			return
		}
		err = userCollection.FindOne(ctx, bson.M{"aadharcardnumber": AadharCard}).Decode(&foundUser)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
			return
		}
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
		tokenString = tokenString[len("Bearer "):]

		err = token.VerifyToken(tokenString)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}

		adminCheck := CheckAdminRole(foundUser)
		if adminCheck {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "admin is not allowed to vote"})
			return
		}

		if foundUser.IsVoted == true {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user has already voted"})
			return
		}

		result, err := userCollection.UpdateOne(
			ctx,
			bson.M{"_id": foundUser.ID},
			bson.D{{"$set", bson.D{
				{"isvoted", true},
			}}},
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to vote"})
			return
		}

		if result.ModifiedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "vote not updated"})
			return
		}
		if err := c.BindJSON(&candidate); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err = candidateCollection.FindOne(ctx, bson.M{"name": candidate.Name}).Decode(&foundCandidate)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "candidate not found"})
			return
		}

		voteDetails.UserID = foundUser.ID
		voteDetails.VotedAt = time.Now()
		count := foundCandidate.VoteCount
		count = count + 1

		result, err = candidateCollection.UpdateOne(
			ctx,
			bson.M{"_id": foundCandidate.ID},
			bson.D{{"$set", bson.D{
				{"votes", voteDetails},
				{"votecount", count},
			}}},
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to vote"})
			return
		}

		if result.ModifiedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "vote not updated"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "vote updated successfully"})
	}
}

func GetVoteCount() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var candidate models.Candidate
		var foundCandidate models.Candidate

		if err := c.BindJSON(&candidate); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := candidateCollection.FindOne(ctx, bson.M{"name": candidate.Name}).Decode(&foundCandidate)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "candidate not found"})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, gin.H{"votecount": foundCandidate.VoteCount})
	}
}

func GetCandidate() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		result, err := candidateCollection.Find(context.TODO(), bson.M{})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing items"})
			return
		}

		var candidateNames []string
		for result.Next(ctx) {
			var candidate models.Candidate
			if err := result.Decode(&candidate); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to decode candidate"})
				return
			}
			candidateNames = append(candidateNames, candidate.Name)
		}
		c.JSON(http.StatusOK, candidateNames)
	}
}
