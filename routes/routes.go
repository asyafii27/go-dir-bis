package routes

import (
	regionController "mobile-directory-bussines/controllers/api/master/region"
	partnerOwnerController "mobile-directory-bussines/controllers/api/partnerowner"
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
				regionGroup.GET("/provinces/:id", regionController.GetProvinceByID)

				regionGroup.GET("/cities", regionController.GetCities)
				regionGroup.GET("/cities/:id", regionController.GetCityByID)

				regionGroup.GET("/districts", regionController.GetDistricts)
				regionGroup.GET("/districts/:id", regionController.GetDistrictByID)

				regionGroup.GET("/villages", regionController.GetVillages)
				regionGroup.GET("/villages/:id", regionController.GetVillageByID)
			}
		}

		partnerOwner := api.Group("/partner-owner")
		{
			partnerOwner.GET("/partner-owners", partnerOwnerController.GetPartnerOwners)
			partnerOwner.GET("/partner-owners/:id", partnerOwnerController.GetPartnerOwnerByID)
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
