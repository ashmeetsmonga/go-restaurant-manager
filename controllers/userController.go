package controller

import (
	"context"
	"golang-restaurant-management/database"
	"golang-restaurant-management/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")

func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100 * time.Second)
		defer cancel()

		var allUsers []models.User
		results, err := userCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch all users"})
			return
		}

		if err := results.All(ctx, &allUsers); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to parse results to allUsers"})
			return
		}

		c.JSON(http.StatusOK, allUsers)
	}
}

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100 * time.Second)
		defer cancel()

		userId := c.Param("user_id")
		var user models.User

		 err := userCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&user)
		 if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch user"})
			return
		 }

		 c.JSON(http.StatusOK, user)
	}
}

func Register() gin.HandlerFunc {
	return func(c *gin.Context) {
		 
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		
	}
}

func HashPassword(password string) string {

}

func VerifyPassword(userPassword string, providedPassword string) (bool, string) {

}