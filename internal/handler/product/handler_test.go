package product_test

// type HandlerTest struct {
// 	e         *echo.Echo
// 	method    string
// 	service   *product_mocks.MockService
// 	presenter *present.Presenter
// 	handler   *product.Handler
// 	req       *http.Request
// 	rec       *httptest.ResponseRecorder
// 	c         echo.Context
// }

// func (t *HandlerTest) SetParam(name, value string) {
// 	t.c.SetParamNames(name)
// 	t.c.SetParamValues(value)
// }

// func setupHandlerTest(t *testing.T, method string, body io.Reader) *HandlerTest {
// 	t.Helper()

// 	e := echo.New()
// 	e.Validator = extendedvalidator.MustNew()

// 	test := &HandlerTest{}
// 	test.e = e
// 	test.method = method
// 	test.service = product_mocks.NewMockService(t)
// 	test.presenter = present.NewPresenter()
// 	ctx := auth.DevUserContext(context.Background())

// 	test.req = httptest.NewRequest(test.method, "/", body).WithContext(ctx)
// 	if body != nil {
// 		test.req.Header.Set("Content-Type", "application/json")
// 	}

// 	test.rec = httptest.NewRecorder()
// 	test.c = test.e.NewContext(test.req, test.rec)

// 	test.handler = product.NewHandler(test.service, test.presenter)

// 	t.Cleanup(func() {
// 		test.service.AssertExpectations(t)
// 	})

// 	return test
// }

// func testProduct1() (domain.Product, view.ProductView) {
// 	id := uuid.MustParse("00000000-0000-0000-0000-000000000001")

// 	data := domain.Product{
// 		User: &domain.User{
// 			ID:    id,
// 			Name:  "User 1",
// 			Email: "user1@example.com",
// 		},
// 		Active:  true,
// 		Product: 12,
// 	}

// 	productView := view.ProductView{
// 		ID: id.String(),
// 		User: &view.ProductUserView{
// 			ID:    id,
// 			Name:  "User 1",
// 			Email: "user1@example.com",
// 		},
// 		Active:  true,
// 		Product: 12,
// 	}

// 	return data, productView
// }

// func testProductData1() (domain.ProductData, view.ProductDataView) {
// 	id := uuid.MustParse("00000000-0000-0000-0000-000000000001")

// 	data := domain.ProductData{
// 		User: &domain.User{
// 			ID:    id,
// 			Name:  "User 1",
// 			Email: "user1@example.com",
// 		},
// 		Active:            true,
// 		Product:           12,
// 		UsedProduct:       10,
// 		LastActiveRequest: nil,
// 	}

// 	productView := view.ProductDataView{
// 		ID: id.String(),
// 		User: &view.ProductUserView{
// 			ID:    id,
// 			Name:  "User 1",
// 			Email: "user1@example.com",
// 		},
// 		Active:         true,
// 		Product:        12,
// 		UsedProduct:    10,
// 		ProductRequest: nil,
// 	}

// 	return data, productView
// }

// func testProductData2() (domain.ProductData, view.ProductDataView) {
// 	id := uuid.MustParse("00000000-0000-0000-0000-000000000002")

// 	data := domain.ProductData{
// 		User: &domain.User{
// 			ID:    id,
// 			Name:  "User 2",
// 			Email: "user2@example.com",
// 		},
// 		Active:      false,
// 		Product:     20,
// 		UsedProduct: 13,
// 		LastActiveRequest: &domain.ProductRequest{
// 			ID:      3,
// 			Status:  domain.ProductRequestStatusPending,
// 			Product: 30,
// 			Remark:  "Some remark",
// 		},
// 	}

// 	productView := view.ProductDataView{
// 		ID: id.String(),
// 		User: &view.ProductUserView{
// 			ID:    id,
// 			Name:  "User 2",
// 			Email: "user2@example.com",
// 		},
// 		Active:      false,
// 		Product:     20,
// 		UsedProduct: 13,
// 		ProductRequest: &view.ProductRequestView{
// 			ID:      3,
// 			Product: 30,
// 			Status:  "Pending",
// 			Remark:  "Some remark",
// 		},
// 	}

// 	return data, productView
// }
