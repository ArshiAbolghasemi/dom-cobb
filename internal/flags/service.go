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

func (s *Service) ValidateCreateFeatureFlagRequest(c *gin.Context) (*CreateFeatureFlagRequest, *api.APIError) {
	var req CreateFeatureFlagRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, api.BadRequestError("Invalid input format", err.Error())
	}

	flag, err := s.repo.GetFlagByName(req.Name)
	if err != nil {
		return nil, api.InternalServerError("Internal Server Error", err.Error())
	}
	if flag != nil {
		return nil, api.ConflictError("Feature flag already exists", "A feature flag with this name already exists")
	}

	if len(req.FeatureFlagIDDependencies) == 0 {
		return &req, nil
	}

	dependencyFlags, err := s.repo.GetFlagByIds(req.FeatureFlagIDDependencies)
	if err != nil {

		return nil, api.InternalServerError("Internal Server Error", err.Error())
	}
	if len(dependencyFlags) != len(req.FeatureFlagIDDependencies) {
		return nil, api.NotFoundError("Invalid dependency feature flag ids", "")
	}

	if req.IsActive {
		if !s.canActivateFlag(dependencyFlags) {
			return nil, api.BadRequestError(
				"Dependency validation failed",
				"Cannot activate feature flag when dependencies are inactive",
			)
		}
	}

	return &req, nil
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
	*FeatureFlag,
	*UpdateFeatureFlagRequest,
	*api.APIError,
) {
	flagId, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return nil, nil, api.BadRequestError("Invalid input format", err.Error())
	}
	var req UpdateFeatureFlagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, nil, api.BadRequestError("Invalid input format", err.Error())
	}
	flag, err := s.repo.GetFlagById(uint(flagId))
	if err != nil {
		return nil, nil, api.InternalServerError("Internal Server Error", err.Error())
	}
	if flag == nil {
		return nil, nil, api.NotFoundError("Invalid flag id", "")
	}
	if req.IsActive == flag.IsActive {
		status := "active"
		if !flag.IsActive {
			status = "inactive"
		}
		return nil, nil, api.OKError(fmt.Sprintf("Flag is already %s", status), "")
	}
	if req.IsActive {
		flagDependencies, err := s.repo.GetFlagDependencies(flag)
		if err != nil {
			return nil, nil, api.InternalServerError("Internal Server Error", err.Error())
		}
		if !s.canActivateFlag(flagDependencies) {
			return nil, nil, api.BadRequestError(
				"Dependencies not satisfied",
				"Cannot activate feature flag when dependencies are inactive",
			)
		}
	} else {
		flagDependents, err := s.repo.GetFlagDependents(flag)
		if err != nil {
			return nil, nil, api.InternalServerError("Internal Server Error", err.Error())
		}
		if !s.canDeactivateFlag(flagDependents) {
			return nil, nil, api.BadRequestError(
				"Dependents still active",
				"Cannot dectivate feature flag when dependents are active",
			)
		}
	}

	return flag, &req, nil
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

func (s *Service) ValidateGetFeatureFlagRequest(c *gin.Context) (*FeatureFlag, *api.APIError) {
	flagId, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return nil, api.BadRequestError("Invalid input format", err.Error())
	}

	flag, err := s.repo.GetFlagById(uint(flagId))
	if err != nil {
		return nil, api.InternalServerError("Internal Server Error", err.Error())
	}
	if flag == nil {
		return nil, api.NotFoundError("Invalid flag id", "")
	}

	return flag, nil
}

func (s *Service) GetFeatureFlag(flag *FeatureFlag) (*FeatureFlagData, error) {
	dependencies, err := s.repo.GetFlagDependencies(flag)
	if err != nil {
		return nil, err
	}

	dependents, err := s.repo.GetFlagDependents(flag)
	if err != nil {
		return nil, err
	}

	dependencyIDs := make([]uint, 0, len(dependencies))
	for _, dependency := range dependencies {
		dependencyIDs = append(dependencyIDs, dependency.ID)
	}

	dependentIDs := make([]uint, 0, len(dependents))
	for _, depentent := range dependents {
		dependentIDs = append(dependentIDs, depentent.ID)
	}

	return &FeatureFlagData{
		ID:           flag.ID,
		Name:         flag.Name,
		Active:       flag.IsActive,
		Dependencies: dependencyIDs,
		Dependents:   dependentIDs,
	}, nil
}
