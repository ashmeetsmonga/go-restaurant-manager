package controller

import (
	"context"
	"golang-restaurant-management/database"
	helper "golang-restaurant-management/helpers"
	"golang-restaurant-management/models"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
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
		ctx, cancel := context.WithTimeout(context.Background(), 100 * time.Second)
		defer cancel()

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

		emailCount, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email});
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to count email"})
			return
		}

		if emailCount > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email already in use"})
		}

		phoneCount, err := userCollection.CountDocuments(ctx, bson.M{"phone": user.Phone});
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to count phone"})
			return
		}

		if phoneCount > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Phone already in use"})
			return
		}

		password := HashPassword(*user.Password)
		user.Password = &password

		user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.User_id = user.ID.Hex()

		token, refreshToken, _ = helper.GenerateAllTokens(*user.Email, *user.First_name, *user.Last_name, user.User_id)
		user.Token = &token
		user.Refresh_Token = &refreshToken
		
		_, err = userCollection.InsertOne(ctx, user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to add user"})
			return
		}

		c.JSON(http.StatusCreated, user)
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		
	}
}

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}

	return string(bytes)
}

func VerifyPassword(userPassword string, providedPassword string) (bool, string) {

}