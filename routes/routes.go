package routes

import (
	controllers "mobile-directory-bussines/controllers/api/user"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	mobile := router.Group("/api/mobile")
	{
		users := mobile.Group("/users")
		{
			users.GET("", controllers.GetUsers)
			users.GET("/:id", controllers.GetUserByID)
			users.POST("", controllers.StoreUser)
			users.PUT("/:id", controllers.UpdateUser)
		}
	}

}
