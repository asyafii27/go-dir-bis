package main

import (
	"mobile-directory-bussines/config"
	"mobile-directory-bussines/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectDatabase()

	router := gin.Default()

	routes.SetupRoutes(router)

	router.Run(":8080")
}
