package flags

import (
	"net/http"

	"github.com/ArshiAbolghasemi/dom-cobb/internal/api"
	"github.com/ArshiAbolghasemi/dom-cobb/internal/logger"
	"github.com/gin-gonic/gin"
)

type CreateFeatureFlagRequest struct {
	Name                      string `json:"name" binding:"required,min=1,max=255"`
	IsActive                  bool   `json:"active"`
	FeatureFlagIDDependencies []uint `json:"feature_flag_id_dependencies"`
}

func CreateFeatureFlagAPI(c *gin.Context) {
	repo := GetRepository()
	logger := logger.NewService()
	service := GetService(repo, logger)

	valid, req := service.ValidateCreateFeatureFlagRequest(c)
	if !valid {
		return
	}

	err := service.CreateFeatureFlag(req)

	if err != nil {
		c.JSON(http.StatusInternalServerError, api.ErrorResponse{
			Error:   "Internal Server Error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, api.SuccessResponse{
		Message: "Feature falg is created successfully",
	})
}
