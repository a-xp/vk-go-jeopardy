package rating

import (
	"github.com/gin-gonic/gin"
	"goj/configuration"
)

var vkKey string

func ConfigureAPI(r *gin.Engine, config configuration.Configuration) {

	vkKey = config.VkApp.Secret

	api := r.Group("/api")
	admin := api.Group("/admin")
	admin.Use()

	api.GET("/me", meEndpoint)
	api.GET("/rating", ratingEndpoint)

	admin.GET("/games", listGamesEndpoint)
	admin.DELETE("/games/:id", removeGameEndpoint)
	admin.POST("/games", updateGameEndpoint)

	admin.GET("/admins", listAdminsEndpoint)
	admin.DELETE("/admins/:id", removeAdminEndpoint)
	admin.POST("/admins", addAdminEndpoint)

	admin.GET("/groups", listGroupsEndpoint)
	admin.DELETE("/groups/:id", removeGroupEndpoint)
	admin.POST("/groups", addGroupEndpoint)

}
