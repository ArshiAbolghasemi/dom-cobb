package flags

import (
	"net/http"

	"github.com/ArshiAbolghasemi/dom-cobb/internal/api"
	"github.com/gin-gonic/gin"
)

type CreateFeatureFlagRequest struct {
	Name                      string `json:"name" binding:"required,min=1,max=255"`
	IsActive                  bool   `json:"active" binding:"required"`
	FeatureFlagIDDependencies []uint `json:"feature_flag_id_dependencies"`
}

func CreateFeatureFlagAPI(c *gin.Context) {
	valid, req, dependencyFlags := ValidateCreateFeatureFlagRequest(c)
	if !valid {
		return
	}

	err := CreateFeatureFlag(req, dependencyFlags)

	if err != nil {
		c.JSON(http.StatusInternalServerError, api.ErrorResponse{
			Error: "Internal Server Error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, api.SuccessResponse{
		Message: "Feature falg is created successfully",
	})
}
