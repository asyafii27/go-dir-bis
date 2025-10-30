package routes

import (
	regionController "mobile-directory-bussines/controllers/api/master/region"
	userController "mobile-directory-bussines/controllers/api/user"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	api := router.Group("/api")
	{
		master := api.Group("/master")
		{
			regionGroup := master.Group("/region")
			{
				regionGroup.GET("/provinces", regionController.GetProvinces)
			}
		}

		mobile := api.Group("/mobile")
		{
			mobile.GET("/users", userController.GetUsers)
			mobile.GET("/users/:id", userController.GetUserByID)
			mobile.POST("/users", userController.StoreUser)
			mobile.PUT("/users/:id", userController.UpdateUser)
		}
	}
}
