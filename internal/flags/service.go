package flags

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/ArshiAbolghasemi/dom-cobb/internal/api"
	"github.com/ArshiAbolghasemi/dom-cobb/internal/logger"
	"github.com/gin-gonic/gin"
)

type Service struct {
	repo   IRepository
	logger logger.IService
}

var (
	service     *Service
	onceService sync.Once
)

func GetService(repo IRepository, logger logger.IService) *Service {
	onceService.Do(func() {
		service = &Service{
			repo:   repo,
			logger: logger,
		}
	})
	return service
}

func (s *Service) ValidateCreateFeatureFlagRequest(c *gin.Context) (*api.APIError, *CreateFeatureFlagRequest) {
	var req CreateFeatureFlagRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		return api.BadRequestError("Invalid input format", err.Error()), nil
	}

	flag, err := s.repo.GetFlagByName(req.Name)
	if err != nil {
		return api.InternalServerError("Internal Server Error", err.Error()), nil
	}
	if flag != nil {
		return api.ConflictError("Feature flag already exists", "A feature flag with this name already exists"), nil
	}

	if len(req.FeatureFlagIDDependencies) == 0 {
		return nil, &req
	}

	dependencyFlags, err := s.repo.GetFlagByIds(req.FeatureFlagIDDependencies)
	if err != nil {

		return api.InternalServerError("Internal Server Error", err.Error()), nil
	}
	if len(dependencyFlags) != len(req.FeatureFlagIDDependencies) {
		return api.NotFoundError("Invalid dependency feature flag ids", ""), nil
	}

	if req.IsActive {
		if !s.canActivateFlag(dependencyFlags) {
			return api.BadRequestError(
				"Dependency validation failed",
				"Cannot activate feature flag when dependencies are inactive",
			), nil
		}
	}

	return nil, &req
}

func (s *Service) CreateFeatureFlag(req *CreateFeatureFlagRequest) error {
	flag, err := s.repo.CreateFlag(req.Name, req.IsActive, req.FeatureFlagIDDependencies)
	if err != nil {
		return err
	}

	s.logger.Log(
		"Feature Flag is created successfully",
		map[string]any{
			"flag_id": flag.ID,
		},
	)

	return nil
}

func (s *Service) ValidateUpdateFeatureFlagRequest(
	c *gin.Context,
) (
	*api.APIError,
	*FeatureFlag,
	*UpdateFeatureFlagRequest,
) {
	flagId, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return api.BadRequestError("Invalid input format", err.Error()), nil, nil
	}
	var req UpdateFeatureFlagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return api.BadRequestError("Invalid input format", err.Error()), nil, nil
	}
	flag, err := s.repo.GetFlagById(uint(flagId))
	if err != nil {
		return api.InternalServerError("Internal Server Error", err.Error()), nil, nil
	}
	if flag == nil {
		return api.NotFoundError("Invalid flag id", ""), nil, nil
	}
	if req.IsActive == flag.IsActive {
		status := "active"
		if !flag.IsActive {
			status = "inactive"
		}
		return api.OKError(fmt.Sprintf("Flag is already %s", status), ""), nil, nil
	}
	if req.IsActive {
		flagDependencies, err := s.repo.GetFlagDependencies(flag)
		if err != nil {
			return api.InternalServerError("Internal Server Error", err.Error()), nil, nil
		}
		if !s.canActivateFlag(flagDependencies) {
			return api.BadRequestError(
				"Dependencies not satisfied",
				"Cannot activate feature flag when dependencies are inactive",
			), nil, nil
		}
	} else {
		flagDependents, err := s.repo.GetFlagDependents(flag)
		if err != nil {
			return api.InternalServerError("Internal Server Error", err.Error()), nil, nil
		}
		if !s.canDeactivateFlag(flagDependents) {
			return api.BadRequestError(
				"Dependents still active",
				"Cannot dectivate feature flag when dependents are active",
			), nil, nil
		}
	}

	return nil, flag, &req
}

func (s *Service) canActivateFlag(flagDependencies []*FeatureFlag) bool {
	for _, depFlag := range flagDependencies {
		if !depFlag.IsActive {
			return false
		}
	}
	return true
}

func (s *Service) canDeactivateFlag(flagDependents []*FeatureFlag) bool {
	for _, depFlag := range flagDependents {
		if depFlag.IsActive {
			return false
		}
	}
	return true
}

func (s *Service) UpdateFeatureFlag(flag *FeatureFlag, req *UpdateFeatureFlagRequest) error {
	err := s.repo.UpdateFlag(flag, req.IsActive)
	if err != nil {
		return err
	}

	s.logger.Log(
		"Feature Flag is toggled successfully",
		map[string]any{
			"flag_id": flag.ID,
			"active":  flag.IsActive,
			"reason":  req.Reason,
		},
	)

	return nil
}
