package flags

import (
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	{
		v1 := router.Group("/api/v1")
		v1.POST("/flags", CreateFeatureFlagAPI)
	}
}
