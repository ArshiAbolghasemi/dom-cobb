package flags

import (
	"net/http"

	"github.com/ArshiAbolghasemi/dom-cobb/internal/api"
	"github.com/ArshiAbolghasemi/dom-cobb/internal/logger"
	"github.com/gin-gonic/gin"
)

func newFeatureFlagService() *Service {
	repo := GetRepository()
	logger := logger.NewService()
	return GetService(repo, logger)
}

// @Description Request payload for creating a new feature flag
type CreateFeatureFlagRequest struct {
	Name                      string `json:"name" binding:"required,min=1,max=255"`
	IsActive                  bool   `json:"active"`
	FeatureFlagIDDependencies []uint `json:"feature_flag_id_dependencies"`
}

// @Summary Create a new feature flag
// @Description Creates a new feature flag with the provided configuration
// @Tags feature-flags
// @Accept json
// @Produce json
// @Param request body CreateFeatureFlagRequest true "Feature flag creation request"
// @Success 201 {object} api.SuccessResponse"Feature flag created successfully"
// @Failure 400 {object} api.ErrorResponse "Bad request - validation error"
// @Failure 404 {object} api.ErrorResponse "Not Found Error"
// @Failure 409 {object} api.ErrorResponse "Conflict Error"
// @Failure 500 {object} api.ErrorResponse "Internal server error"
// @Router /api/v1/flags [post]
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

// @Description Request payload for updating a feature flag
type UpdateFeatureFlagRequest struct {
	IsActive bool   `json:"active"`
	Reason   string `json:"reason" binding:"required,min=1,max=255"`
}

// @Summary Update a feature flag
// @Description Update a feature flag with the provided configuration
// @Tags feature-flags
// @Accept json
// @Produce json
// @Param request body UpdateFeatureFlagRequest true "Feature flag creation request"
// @Success 200 {object} api.SuccessResponse"Feature flag is updated successfully"
// @Failure 400 {object} api.ErrorResponse "Bad request - validation error"
// @Failure 404 {object} api.ErrorResponse "Not Found Error"
// @Failure 500 {object} api.ErrorResponse "Internal server error"
// @Router /api/v1/flags/{id} [patch]
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

// @Description Feature flag data with dependencies and dependents information
type FeatureFlagData struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	Active       bool   `json:"active"`
	Dependencies []uint `json:"dependencies"`
	Dependents   []uint `json:"dependents"`
}

// @Summary Get a feature flag
// @Description Retrieve a feature flag by its ID including its dependencies and dependents
// @Tags feature-flags
// @Accept json
// @Produce json
// @Param id path int true "Feature Flag ID"
// @Success 200 {object} api.SuccessResponse{data=FeatureFlagData} "Feature flag retrieved successfully"
// @Failure 400 {object} api.ErrorResponse "Bad request - validation error"
// @Failure 404 {object} api.ErrorResponse "Feature flag not found"
// @Failure 500 {object} api.ErrorResponse "Internal server error"
// @Router /api/v1/flags/{id} [get]
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

// @Description Query parameters for paginated feature flag logs request
type GetFeatureFlagLogsQueryParams struct {
	api.PaginationQueryParam
}

// @Description Paginated response containing feature flag logs
type GetFeatureFlagLogsData struct {
	Logs []*logger.LogEntry `json:"logs"`
	api.PaginationResponse
}

// @Summary Get feature flag logs
// @Description Retrieve paginated logs for a specific feature flag
// @Tags feature-flags
// @Accept json
// @Produce json
// @Param id path int true "Feature Flag ID"
// @Param page query int false "Page number (default: 1)" minimum(1)
// @Param size query int false "Number of items per page (default: 10)" minimum(1) maximum(20)
// @Success 200 {object} api.SuccessResponse{data=GetFeatureFlagLogsData} "Feature flag logs retrieved successfully"
// @Failure 400 {object} api.ErrorResponse "Bad request - validation error"
// @Failure 404 {object} api.ErrorResponse "Feature flag not found"
// @Failure 500 {object} api.ErrorResponse "Internal server error"
// @Router /api/v1/flags/{id}/logs [get]
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
