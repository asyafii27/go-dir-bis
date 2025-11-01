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
				// privinces
				regionGroup.GET("/provinces", regionController.GetProvinces)
				regionGroup.GET("/provinces/:id", regionController.GetProvinceByID)

				// cities
				regionGroup.GET("cities", regionController.GetCities)
				regionGroup.GET("cities/:id", regionController.GetCityByID)
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
