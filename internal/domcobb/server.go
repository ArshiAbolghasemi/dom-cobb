package domcobb

import "github.com/gin-gonic/gin"

func Run() {
	r := gin.Default()

	SetupRoutes(r)

	port, err := GetAppPort()
	if err != nil {
		panic(err)
	}
	r.Run(":" + port)
}
