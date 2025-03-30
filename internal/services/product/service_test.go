package product_test

// import (
// 	"salehelp/internal/pkg/logger/handlers/slogdiscard"
// 	"salehelp/internal/service/product"
// 	"salehelp/mocks/salehelp/int/service/product_mocks"
// 	"testing"
// )

// type ServiceTest struct {
// 	srv          *product.Service
// 	repo         *product_mocks.MockRepository
// 	requestRepo  *product_mocks.MockRequestRepository
// 	appRepo      *product_mocks.MockApplicationRepository
// 	eventService *product_mocks.MockEventService
// }

// func setupServiceTest(t *testing.T) *ServiceTest {
// 	t.Helper()

// 	repoMock := product_mocks.NewMockRepository(t)
// 	requestRepoMock := product_mocks.NewMockRequestRepository(t)
// 	appRepoMock := product_mocks.NewMockApplicationRepository(t)
// 	eventServiceMock := product_mocks.NewMockEventService(t)

// 	service := product.NewService(slogdiscard.NewDiscardLogger(), repoMock, requestRepoMock, appRepoMock, eventServiceMock)

// 	t.Cleanup(func() {
// 		repoMock.AssertExpectations(t)
// 		requestRepoMock.AssertExpectations(t)
// 		appRepoMock.AssertExpectations(t)
// 		eventServiceMock.AssertExpectations(t)
// 	})

// 	return &ServiceTest{
// 		srv:          service,
// 		repo:         repoMock,
// 		requestRepo:  requestRepoMock,
// 		appRepo:      appRepoMock,
// 		eventService: eventServiceMock,
// 	}
// }
