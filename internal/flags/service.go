package flags

import (
	"net/http"

	"github.com/ArshiAbolghasemi/dom-cobb/internal/api"
	"github.com/ArshiAbolghasemi/dom-cobb/internal/database/postgres"
	"github.com/gin-gonic/gin"
)

func ValidateCreateFeatureFlagRequest(c *gin.Context) (bool, *CreateFeatureFlagRequest, []FeatureFlag) {
	var req CreateFeatureFlagRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, api.ErrorResponse{
			Error:   "Invalid input format",
			Message: err.Error(),
		})
		return false, nil, nil
	}

	flag, err := GetFlagByName(req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, api.ErrorResponse{
			Error:   "Internal Server Error",
			Message: err.Error(),
		})
		return false, nil, nil
	}
	if flag != nil {
		c.JSON(http.StatusConflict, api.ErrorResponse{
			Error:   "Feature flag already exists",
			Message: "A feature flag with this name already exists",
		})
		return false, nil, nil
	}

	if len(req.FeatureFlagIDDependencies) > 0 {
		dependencyFlags, err := GetFlagByIds(req.FeatureFlagIDDependencies)
		if err != nil {
			c.JSON(http.StatusInternalServerError, api.ErrorResponse{
				Error:   "Internal Server Error",
				Message: err.Error(),
			})
			return false, nil, nil
		}
		if len(dependencyFlags) != len(req.FeatureFlagIDDependencies) {
			c.JSON(http.StatusBadRequest, api.ErrorResponse{
				Error: "Invalid dependency feature flag ids",
			})
			return false, nil, nil
		}

		if req.IsActive {
			for _, depFlag := range dependencyFlags {
				if !depFlag.IsActive {
					c.JSON(http.StatusBadRequest, api.ErrorResponse{
						Error:   "Dependency validation failed",
						Message: "Cannot activate feature flag when dependency '" + depFlag.Name + "' is inactive",
					})
					return false, nil, nil
				}
			}
		}

		return true, &req, dependencyFlags
	} else {
		return true, &req, nil
	}
}

func CreateFeatureFlag(req *CreateFeatureFlagRequest, depenedencyFlags []FeatureFlag) error {
	flag := FeatureFlag{
		Name:     req.Name,
		IsActive: req.IsActive,
	}
	db := postgres.GetDB()

	if len(depenedencyFlags) == 0 {
		err := db.Create(&flag).Error
		return err
	}

	tx := db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	if err := tx.Create(&flag).Error; err != nil {
		tx.Rollback()
		return err
	}

	var dependencies []FlagDependency
	for _, depFlag := range depenedencyFlags {
		dependencies = append(dependencies, FlagDependency{
			FlagID:          flag.ID,
			DependsOnFlagID: depFlag.ID,
		})
	}

	if err := tx.Create(&dependencies).Error; err != nil {
		tx.Rollback()
		return err
	}

	err := tx.Commit().Error
	return err
}
