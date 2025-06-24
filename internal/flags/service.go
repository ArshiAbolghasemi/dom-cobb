package flags

import (
	"net/http"
	"sync"

	"github.com/ArshiAbolghasemi/dom-cobb/internal/api"
	"github.com/gin-gonic/gin"
)

type Service struct {
	repo IRepository
}

var (
	service     *Service
	onceService sync.Once
)

func GetService(repo IRepository) *Service {
	onceService.Do(func() {
		service = &Service{
			repo: repo,
		}
	})
	return service
}

func (s *Service) ValidateCreateFeatureFlagRequest(c *gin.Context) (bool, *CreateFeatureFlagRequest) {
	var req CreateFeatureFlagRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, api.ErrorResponse{
			Error:   "Invalid input format",
			Message: err.Error(),
		})
		return false, nil
	}

	flag, err := s.repo.GetFlagByName(req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, api.ErrorResponse{
			Error:   "Internal Server Error",
			Message: err.Error(),
		})
		return false, nil
	}
	if flag != nil {
		c.JSON(http.StatusConflict, api.ErrorResponse{
			Error:   "Feature flag already exists",
			Message: "A feature flag with this name already exists",
		})
		return false, nil
	}

	if len(req.FeatureFlagIDDependencies) == 0 {
		return true, &req
	}

	dependencyFlags, err := s.repo.GetFlagByIds(req.FeatureFlagIDDependencies)
	if err != nil {
		c.JSON(http.StatusInternalServerError, api.ErrorResponse{
			Error:   "Internal Server Error",
			Message: err.Error(),
		})
		return false, nil
	}
	if len(dependencyFlags) != len(req.FeatureFlagIDDependencies) {
		c.JSON(http.StatusBadRequest, api.ErrorResponse{
			Error: "Invalid dependency feature flag ids",
		})
		return false, nil
	}

	if req.IsActive {
		for _, depFlag := range dependencyFlags {
			if !depFlag.IsActive {
				c.JSON(http.StatusBadRequest, api.ErrorResponse{
					Error:   "Dependency validation failed",
					Message: "Cannot activate feature flag when dependency '" + depFlag.Name + "' is inactive",
				})
				return false, nil
			}
		}
	}

	return true, &req
}

func (s *Service) CreateFeatureFlag(req *CreateFeatureFlagRequest) error {
	flag := FeatureFlag{
		Name:     req.Name,
		IsActive: req.IsActive,
	}

	var dependencies []FlagDependency
	for _, depFlagID := range req.FeatureFlagIDDependencies {
		dependencies = append(dependencies, FlagDependency{
			FlagID:          flag.ID,
			DependsOnFlagID: depFlagID,
		})
	}

	return s.repo.CreateFlag(flag, dependencies)
}
