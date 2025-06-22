package main

import (
	"github.com/ArshiAbolghasemi/dom-cobb/internal/domcobb"
	"github.com/ArshiAbolghasemi/dom-cobb/internal/flags"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	flags.SetupRoutes(r)
	port, err := domcobb.GetAppPort()
	if err != nil {
		panic(err)
	}
	r.Run(":" + port)
}
