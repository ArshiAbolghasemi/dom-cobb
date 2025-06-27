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

func newFeatureFlagService() *Service {
	repo := GetRepository()
	logger := logger.NewService()
	return GetService(repo, logger)
}

func CreateFeatureFlagAPI(c *gin.Context) {
	service := newFeatureFlagService()

	req, err := service.ValidateCreateFeatureFlagRequest(c)
	if err != nil {
		api.RespondAPIError(c, err)
		return
	}

	err = service.CreateFeatureFlag(req)
	if err != nil {
		api.RespondAPIError(c, err)
		return
	}

	api.RespondSuccess(c, http.StatusCreated, "Feature flag is created successfully", nil)
}

type UpdateFeatureFlagRequest struct {
	IsActive bool   `json:"active"`
	Reason   string `json:"reason" binding:"required,min=1,max=255"`
}

func UpdateFeatureFlagAPI(c *gin.Context) {
	service := newFeatureFlagService()

	flag, req, err := service.ValidateUpdateFeatureFlagRequest(c)
	if err != nil {
		api.RespondAPIError(c, err)
		return
	}

	err = service.UpdateFeatureFlag(flag, req)
	if err != nil {
		api.RespondAPIError(c, err)
		return
	}

	api.RespondSuccess(c, http.StatusOK, "Feature flag is updated successfully", nil)
}

type FeatureFlagData struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	Active       bool   `json:"active"`
	Dependencies []uint `json:"dependencies"`
	Dependents   []uint `json:"dependents"`
}

func GetFeatureFlagAPI(c *gin.Context) {
	service := newFeatureFlagService()

	flag, err := service.ValidateGetFeatureFlagRequest(c)
	if err != nil {
		api.RespondAPIError(c, err)
		return
	}

	data, err := service.GetFeatureFlag(flag)
	if err != nil {
		api.RespondAPIError(c, err)
		return
	}

	api.RespondSuccess(c, http.StatusOK, "Feature flag is retrieved successfully", data)
}

type GetFeatureFlagLogsQueryParams struct {
	api.PaginationQueryParam
}

type GetFeatureFlagLogsData struct {
	Logs []*logger.LogEntry `json:"logs"`
	api.PaginationResponse
}

func GetFeatureFlagLogsAPI(c *gin.Context) {
	service := newFeatureFlagService()

	query, flag, apiErr := service.ValidateGetFeatureFlagLogsRequest(c)
	if apiErr != nil {
		api.RespondAPIError(c, apiErr)
		return
	}

	data, err := service.GetFeatureFlagLogs(flag, query)
	if err != nil {
		api.RespondAPIError(c, err)
		return
	}

	api.RespondSuccess(c, http.StatusOK, "Feature flag logs is retrieved successfully", data)
}
