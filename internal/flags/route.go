package flags

import (
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	{
		v1 := router.Group("/api/v1")
		v1.POST("/flags", CreateFeatureFlagAPI)
		v1.PATCH("/flags/:id", UpdateFeatureFlagAPI)
		v1.GET("/flags/:id", GetFeatureFlagAPI)
		v1.GET("/flags/:id/logs", GetFeatureFlagLogsAPI)
	}
}
