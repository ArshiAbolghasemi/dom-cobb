package main

import (
	_ "github.com/ArshiAbolghasemi/dom-cobb/docs"
	"github.com/ArshiAbolghasemi/dom-cobb/internal/domcobb"
)

// @title           Dom Cobb API
// @version         1.0
// @description     This is a Dom Cobb API documentation

// @host      localhost:8080
// @BasePath  /api/v1
func main() {
	domcobb.Run()
}
