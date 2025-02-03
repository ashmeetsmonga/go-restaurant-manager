package main

import (
	"fmt"
	"golang-restaurant-management/middleware"
	"golang-restaurant-management/routes"
	"os"

	"github.com/gin-gonic/gin"
)


func main() {
	port := os.Getenv("PORT")
	if port == "" {
		fmt.Println("Unable to load from env file")
		port = "8080"
	}

	router := gin.New()
	router.Use(gin.Logger())
	routes.UserRoutes(router)
	router.Use(middleware.Authentication())

	routes.FoodRoutes(router)
	routes.InvoiceRoutes(router)
	routes.MenuRoutes(router)
	routes.OrderItemRoutes(router)
	routes.OrderRoutes(router)
	routes.TableRoutes(router)

	router.Run(":" + port)
}