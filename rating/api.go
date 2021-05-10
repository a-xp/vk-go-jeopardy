package rating

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"goj/configuration"
	"time"
)

var vkKey string
var validateSignature bool

func ConfigureAPI(r *gin.Engine, config *configuration.Configuration) {

	vkKey = config.VkApp.Secret
	validateSignature = config.ValidateRequest

	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"POST", "GET", "DELETE", "PUT"},
		AllowHeaders:     []string{"X-VK-PARAMS", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
		MaxAge:           24 * time.Hour,
	}))

	api := r.Group("/api")
	admin := api.Group("/admin")
	api.Use(FilterPublic)
	admin.Use(FilterPublic, FilterAdmin)

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
