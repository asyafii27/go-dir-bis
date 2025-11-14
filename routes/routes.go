package routes

import (
	categoryController "mobile-directory-bussines/controllers/api/master/category"
	regionController "mobile-directory-bussines/controllers/api/master/region"
	secondSubCategoryController "mobile-directory-bussines/controllers/api/master/secondsubcategory"
	subCategoryController "mobile-directory-bussines/controllers/api/master/subcategory"
	partnerController "mobile-directory-bussines/controllers/api/partner"
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

			// categories
			master.GET("/categories", categoryController.GetCategories)
			master.GET("/categories/:id", categoryController.GetCategoryByID)
			master.POST("/categories", categoryController.StoreCategory)

			// sub categories
			master.GET("/subcategories", subCategoryController.GetSubCategories)
			master.GET("/subcategories/:id", subCategoryController.GetSubCategoryByID)

			// second sub category
			master.GET("/secondsubcategories", secondSubCategoryController.GetSecondSubCategories)
			master.GET("/secondsubcategories/:id", secondSubCategoryController.GetSecondSubCategoryByID)
		}

		// partner owners
		partnerOwner := api.Group("/partner-owner")
		{
			partnerOwner.GET("/partner-owners", partnerOwnerController.GetPartnerOwners)
			partnerOwner.GET("/partner-owners/:id", partnerOwnerController.GetPartnerOwnerByID)
		}

		// partners
		partner := api.Group("/partner")
		{
			partner.GET("/partners", partnerController.GetPartners)
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
