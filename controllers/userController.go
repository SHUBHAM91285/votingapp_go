package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/SHUBHAM91285/votingApp_go/database"
	"github.com/SHUBHAM91285/votingApp_go/models"
	"github.com/SHUBHAM91285/votingApp_go/token"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")

func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		adminCount, err := userCollection.CountDocuments(ctx, bson.M{"role": "admin"})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking admin count"})
			return
		}

		if adminCount > 0 && user.Role == "admin" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Only one admin user is allowed"})
			return
		}

		err = userCollection.FindOne(ctx, bson.M{"aadharcardnumber": user.AadharCardNumber}).Decode(&user)
		if err == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user already exist"})
			return
		}
		user.ID = primitive.NewObjectID()
		password := HashPassword(user.Password)
		user.Password = password
		tokenString, err := token.CreateToken(user.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}

		userInfo, insertErr := userCollection.InsertOne(ctx, user)
		if insertErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User is not created"})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, gin.H{
			"userInfo": userInfo,
			"token":    tokenString,
		})
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User
		var foundUser models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		err := userCollection.FindOne(ctx, bson.M{"aadharcardnumber": user.AadharCardNumber}).Decode(&foundUser)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found,login seems to be incorrect"})
			return
		}
		passwordIsValid := VerifyPassword(user.Password, foundUser.Password)
		defer cancel()
		if passwordIsValid != true {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "password is invalid"})
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

		c.JSON(http.StatusOK, foundUser)
	}
}

func UserProfile() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
		tokenString = tokenString[len("Bearer "):]

		err := token.VerifyToken(tokenString)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}

		var user models.User
		var foundUser models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		err = userCollection.FindOne(ctx, bson.M{"aadharcardnumber": user.AadharCardNumber}).Decode(&foundUser)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found,login seems to be incorrect"})
			return
		}

		c.JSON(http.StatusOK, foundUser)
	}
}

func UpdatePassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
		tokenString = tokenString[len("Bearer "):]

		err := token.VerifyToken(tokenString)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}

		var user models.User
		var foundUser models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err = userCollection.FindOne(ctx, bson.M{"aadharcardnumber": user.AadharCardNumber}).Decode(&foundUser)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found,login seems to be incorrect"})
			return
		}

		hashedPassword := HashPassword(user.Password)
		update := bson.M{"$set": bson.M{"password": hashedPassword}}
		_, err = userCollection.UpdateOne(ctx, bson.M{"aadharcardnumber": user.AadharCardNumber}, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error updating password"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "password updated successfully"})

	}
}

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}

func VerifyPassword(userPassword, providedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	check := true
	if err != nil {
		check = false
	}
	return check
}
