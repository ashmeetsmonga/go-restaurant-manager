package middleware

import (
	helper "golang-restaurant-management/helpers"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No valid Authorization header provided"})
			c.Abort()
			return
		}
		clientToken := strings.TrimPrefix(authHeader, "Bearer ")
		if clientToken == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No Auth Token provided"})
			c.Abort()
			return
		}

		claims, err := helper.ValidateToken(clientToken)
		if err != "" {
			c.JSON(http.StatusForbidden, gin.H{"error": err})
			c.Abort()
			return
		}

		c.Set("email", claims.Email)
		c.Set("first_name", claims.First_name)
		c.Set("last_name", claims.Last_name)
		c.Set("uid", claims.Uid)

		c.Next()
	}
}