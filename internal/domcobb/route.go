package domcobb

import (
	"github.com/ArshiAbolghasemi/dom-cobb/internal/flags"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	flags.SetupRoutes(router)
}
