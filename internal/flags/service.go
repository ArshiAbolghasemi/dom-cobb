package flags

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/ArshiAbolghasemi/dom-cobb/internal/api"
	"github.com/ArshiAbolghasemi/dom-cobb/internal/logger"
	"github.com/gin-gonic/gin"
)

type Service struct {
	Repo   IRepository
	Logger logger.IService
}

var (
	service     *Service
	onceService sync.Once
)

func GetService(repo IRepository, logger logger.IService) *Service {
	onceService.Do(func() {
		service = &Service{
			Repo:   repo,
			Logger: logger,
		}
	})
	return service
}

func (s *Service) ValidateCreateFeatureFlagRequest(c *gin.Context) (*CreateFeatureFlagRequest, *api.APIError) {
	var req CreateFeatureFlagRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, api.BadRequestError("Invalid input format", err.Error())
	}

	flag, err := s.Repo.GetFlagByName(req.Name)
	if err != nil {
		return nil, api.InternalServerError("Internal Server Error", err.Error())
	}
	if flag != nil {
		return nil, api.ConflictError("Feature flag already exists", "A feature flag with this name already exists")
	}

	if len(req.FeatureFlagIDDependencies) == 0 {
		return &req, nil
	}

	dependencyFlags, err := s.Repo.GetFlagByIds(req.FeatureFlagIDDependencies)
	if err != nil {
		return nil, api.InternalServerError("Internal Server Error", err.Error())
	}
	if len(dependencyFlags) != len(req.FeatureFlagIDDependencies) {
		return nil, api.NotFoundError("Invalid dependency feature flag ids", "")
	}

	if req.IsActive {
		if canActivate, inactiveIds := s.canActivateFlag(dependencyFlags); !canActivate {
			return nil, api.BadRequestError(
				"Dependency validation failed",
				fmt.Sprintf("Cannot activate feature flag. Missing dependency IDs: %v", inactiveIds),
			)
		}
	}

	return &req, nil
}

func (s *Service) CreateFeatureFlag(req *CreateFeatureFlagRequest) *api.APIError {
	flag, err := s.Repo.CreateFlag(req.Name, req.IsActive, req.FeatureFlagIDDependencies)
	if err != nil {
		return api.InternalServerError("Internal Server Error", err.Error())
	}

	s.Logger.Log(&logger.LogEntry{
		Message: "Feature Flag is created successfully",
		Metadata: map[string]any{
			"flag_id": flag.ID,
		},
		Timestamp: time.Now(),
	})

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
	flag, err := s.Repo.GetFlagById(uint(flagId))
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

	if !req.IsActive {
		return flag, &req, nil
	}

	flagDependencies, err := s.Repo.GetFlagDependencies(flag)
	if err != nil {
		return nil, nil, api.InternalServerError("Internal Server Error", err.Error())
	}
	if canActivate, inactiveIds := s.canActivateFlag(flagDependencies); !canActivate {
		return nil, nil, api.BadRequestError(
			"Dependency validation failed",
			fmt.Sprintf("Cannot activate feature flag. Missing dependency IDs: %v", inactiveIds),
		)
	}

	return flag, &req, nil
}

func (s *Service) canActivateFlag(flagDependencies []*FeatureFlag) (bool, []uint) {
	var inactiveIds []uint
	for _, depFlag := range flagDependencies {
		if !depFlag.IsActive {
			inactiveIds = append(inactiveIds, depFlag.ID)
		}
	}
	return len(inactiveIds) == 0, inactiveIds
}

func (s *Service) UpdateFeatureFlag(flag *FeatureFlag, req *UpdateFeatureFlagRequest) *api.APIError {
	err := s.Repo.UpdateFlag(flag, req.IsActive)
	if err != nil {
		return api.InternalServerError("Internal Server Error", err.Error())
	}

	s.Logger.Log(&logger.LogEntry{
		Message: "Feature Flag is toggled successfully",
		Metadata: map[string]any{
			"flag_id": flag.ID,
			"active":  flag.IsActive,
			"reason":  req.Reason,
		},
		Timestamp: time.Now(),
	})

	return nil
}

func (s *Service) ValidateGetFeatureFlagRequest(c *gin.Context) (*FeatureFlag, *api.APIError) {
	flagId, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return nil, api.BadRequestError("Invalid input format", err.Error())
	}

	flag, err := s.Repo.GetFlagById(uint(flagId))
	if err != nil {
		return nil, api.InternalServerError("Internal Server Error", err.Error())
	}
	if flag == nil {
		return nil, api.NotFoundError("Invalid flag id", "")
	}

	return flag, nil
}

func (s *Service) GetFeatureFlag(flag *FeatureFlag) (*FeatureFlagData, *api.APIError) {
	dependencies, err := s.Repo.GetFlagDependencies(flag)
	if err != nil {
		return nil, api.InternalServerError("Internal Server Error", err.Error())
	}

	dependents, err := s.Repo.GetFlagDependents(flag)
	if err != nil {
		return nil, api.InternalServerError("Internal Server Error", err.Error())
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

func (s *Service) ValidateGetFeatureFlagLogsRequest(
	c *gin.Context,
) (
	*GetFeatureFlagLogsQueryParams,
	*FeatureFlag,
	*api.APIError,
) {
	flagId, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return nil, nil, api.BadRequestError("Invalid input format", err.Error())
	}
	flag, err := s.Repo.GetFlagById(uint(flagId))
	if err != nil {
		return nil, nil, api.InternalServerError("Internal Server Error", err.Error())
	}
	if flag == nil {
		return nil, nil, api.NotFoundError("Invalid flag id", "")
	}

	var query GetFeatureFlagLogsQueryParams
	if err := c.ShouldBindQuery(&query); err != nil {
		return nil, nil, api.BadRequestError("Invalid input format", err.Error())
	}

	return &query, flag, nil
}

func (s *Service) GetFeatureFlagLogs(
	flag *FeatureFlag,
	query *GetFeatureFlagLogsQueryParams,
) (
	*GetFeatureFlagLogsData,
	*api.APIError,
) {
	logs, total, totalPages, err := s.Repo.GetFeatureFlagLogs(flag, query.Page, query.Size)
	if err != nil {
		return nil, api.InternalServerError("Internal Server Error", err.Error())
	}

	return &GetFeatureFlagLogsData{
		Logs: logs,
		PaginationResponse: api.PaginationResponse{
			Page:       query.Page,
			Size:       query.Size,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}
