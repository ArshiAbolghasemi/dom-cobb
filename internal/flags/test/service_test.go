package flags_test

import (
	"net/http"

	"github.com/ArshiAbolghasemi/dom-cobb/internal/flags"
	mockFlags "github.com/ArshiAbolghasemi/dom-cobb/internal/flags/test/mock"
	mockLogger "github.com/ArshiAbolghasemi/dom-cobb/internal/logger/test/mock"
	"github.com/ArshiAbolghasemi/dom-cobb/internal/testutils"
	"github.com/brianvoe/gofakeit/v7"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Service", func() {
	Describe("Validate Create Feature Flag Request", func() {
		When("Request is valid", func() {
			DescribeTable("should return instance create feature flag requests with no error",
				func(
					setupMock func() (*flags.Service, *mockFlags.MockRepository, *mockLogger.MockLogger),
					req *flags.CreateFeatureFlagRequest,
				) {
					service, repo, logger := setupMock()
					c, _ := testutils.CreateJSONRequest(http.MethodPost, "/api/v1/flags", req)

					result, err := service.ValidateCreateFeatureFlagRequest(c)

					Expect(err).To(BeNil())
					Expect(result).NotTo(BeNil())
					Expect(*result).To(Equal(*req))

					repo.AssertExpectations(GinkgoT())
					logger.AssertExpectations(GinkgoT())
				},
				Entry(
					"inactive with no dependency flag",
					func() (*flags.Service, *mockFlags.MockRepository, *mockLogger.MockLogger) {
						r := &mockFlags.MockRepository{}
						r.On("GetFlagByName", "order").Return(nil, nil)
						l := &mockLogger.MockLogger{}
						return &flags.Service{
							Repo:   r,
							Logger: l,
						}, r, l
					},
					&flags.CreateFeatureFlagRequest{
						Name:     "order",
						IsActive: false,
					},
				),
				Entry(
					"inactive with dependency flag",
					func() (*flags.Service, *mockFlags.MockRepository, *mockLogger.MockLogger) {
						r := &mockFlags.MockRepository{}
						r.On("GetFlagByName", "dashboard").Return(nil, nil)
						r.On("GetFlagByIds", []uint{1, 2}).Return(mockFlags.CreateFeatureFlagByIds([]uint{1, 2}), nil)
						l := &mockLogger.MockLogger{}
						return &flags.Service{
							Repo:   r,
							Logger: l,
						}, r, l
					},
					&flags.CreateFeatureFlagRequest{
						Name:                      "dashboard",
						IsActive:                  false,
						FeatureFlagIDDependencies: []uint{1, 2},
					},
				),
				Entry(
					"active with dependency flag",
					func() (*flags.Service, *mockFlags.MockRepository, *mockLogger.MockLogger) {
						r := &mockFlags.MockRepository{}
						r.On("GetFlagByName", "otp").Return(nil, nil)
						r.On("GetFlagByIds", []uint{1, 2, 3}).Return(
							mockFlags.CreateFeatureFlagByIds([]uint{1, 2, 3}, mockFlags.WithIsActive(true)),
							nil,
						)
						l := &mockLogger.MockLogger{}
						return &flags.Service{
							Repo:   r,
							Logger: l,
						}, r, l
					},
					&flags.CreateFeatureFlagRequest{
						Name:                      "otp",
						IsActive:                  true,
						FeatureFlagIDDependencies: []uint{1, 2, 3},
					},
				),
			)
		})

		When("Request is invalida", func() {
			DescribeTable("should return error with empty instance of request",
				func(
					setupMock func() (*flags.Service, *mockFlags.MockRepository, *mockLogger.MockLogger),
					req *flags.CreateFeatureFlagRequest,
					httpCode int,
				) {
					service, repo, logger := setupMock()
					c, _ := testutils.CreateJSONRequest(http.MethodPost, "/api/v1/flags", req)

					result, err := service.ValidateCreateFeatureFlagRequest(c)

					Expect(err).NotTo(BeNil())
					Expect(err.StatusCode).To(Equal(httpCode))
					Expect(result).To(BeNil())

					repo.AssertExpectations(GinkgoT())
					logger.AssertExpectations(GinkgoT())
				},
				Entry(
					"name is missed",
					func() (*flags.Service, *mockFlags.MockRepository, *mockLogger.MockLogger) {
						r := &mockFlags.MockRepository{}
						l := &mockLogger.MockLogger{}
						return &flags.Service{Repo: r, Logger: l}, r, l
					},
					&flags.CreateFeatureFlagRequest{
						IsActive: true,
					},
					http.StatusBadRequest,
				),
				Entry(
					"already flag existed with same name",
					func() (*flags.Service, *mockFlags.MockRepository, *mockLogger.MockLogger) {
						r := &mockFlags.MockRepository{}
						r.On("GetFlagByName", "order").Return(
							mockFlags.CreateFeatureFlag(mockFlags.WithName("order")),
							nil,
						)
						l := &mockLogger.MockLogger{}
						return &flags.Service{
							Repo:   r,
							Logger: l,
						}, r, l
					},
					&flags.CreateFeatureFlagRequest{
						Name:     "order",
						IsActive: false,
					},
					http.StatusConflict,
				),
				Entry(
					"invalid dependency flag ids",
					func() (*flags.Service, *mockFlags.MockRepository, *mockLogger.MockLogger) {
						r := &mockFlags.MockRepository{}
						r.On("GetFlagByName", "dashboard").Return(nil, nil)
						r.On("GetFlagByIds", []uint{1, 2}).Return(
							mockFlags.CreateFeatureFlagByIds([]uint{1}),
							nil,
						)
						l := &mockLogger.MockLogger{}
						return &flags.Service{
							Repo:   r,
							Logger: l,
						}, r, l
					},
					&flags.CreateFeatureFlagRequest{
						Name:                      "dashboard",
						IsActive:                  false,
						FeatureFlagIDDependencies: []uint{1, 2},
					},
					http.StatusNotFound,
				),
				Entry(
					"active flag with inactive dependency",
					func() (*flags.Service, *mockFlags.MockRepository, *mockLogger.MockLogger) {
						r := &mockFlags.MockRepository{}
						r.On("GetFlagByName", "otp").Return(nil, nil)
						r.On("GetFlagByIds", []uint{1}).Return(
							mockFlags.CreateFeatureFlagByIds([]uint{1}, mockFlags.WithIsActive(false)),
							nil,
						)
						l := &mockLogger.MockLogger{}
						return &flags.Service{
							Repo:   r,
							Logger: l,
						}, r, l
					},
					&flags.CreateFeatureFlagRequest{
						Name:                      "otp",
						IsActive:                  true,
						FeatureFlagIDDependencies: []uint{1},
					},
					http.StatusBadRequest,
				),
			)
		})

		When("Internal Server Error is happened", func() {
			It("should return api error with status code 500", func() {
				repo := &mockFlags.MockRepository{}
				repo.On("GetFlagByName", "otp").Return(nil, gofakeit.ErrorDatabase())
				logger := &mockLogger.MockLogger{}
				service := &flags.Service{
					Repo:   repo,
					Logger: logger,
				}

				req := flags.CreateFeatureFlagRequest{
					Name:     "otp",
					IsActive: true,
				}
				c, _ := testutils.CreateJSONRequest(http.MethodPost, "/api/v1/flags", req)

				result, err := service.ValidateCreateFeatureFlagRequest(c)

				Expect(err).NotTo(BeNil())
				Expect(err.StatusCode).To(Equal(http.StatusInternalServerError))
				Expect(result).To(BeNil())

				repo.AssertExpectations(GinkgoT())
				logger.AssertExpectations(GinkgoT())
			})
		})
	})
})
