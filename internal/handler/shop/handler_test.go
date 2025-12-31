package shop_test

// type HandlerTest struct {
// 	e         *echo.Echo
// 	method    string
// 	service   *shop_mocks.MockService
// 	presenter *present.Presenter
// 	handler   *shop.Handler
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
// 	test.service = shop_mocks.NewMockService(t)
// 	test.presenter = present.NewPresenter()
// 	ctx := auth.DevUserContext(context.Background())

// 	test.req = httptest.NewRequest(test.method, "/", body).WithContext(ctx)
// 	if body != nil {
// 		test.req.Header.Set("Content-Type", "application/json")
// 	}

// 	test.rec = httptest.NewRecorder()
// 	test.c = test.e.NewContext(test.req, test.rec)

// 	test.handler = shop.NewHandler(test.service, test.presenter)

// 	t.Cleanup(func() {
// 		test.service.AssertExpectations(t)
// 	})

// 	return test
// }

// func testshop1() (domain.shop, view.shopView) {
// 	id := uuid.MustParse("00000000-0000-0000-0000-000000000001")

// 	data := domain.shop{
// 		User: &domain.User{
// 			ID:    id,
// 			Name:  "User 1",
// 			Email: "user1@example.com",
// 		},
// 		Active:  true,
// 		shop: 12,
// 	}

// 	shopView := view.shopView{
// 		ID: id.String(),
// 		User: &view.shopUserView{
// 			ID:    id,
// 			Name:  "User 1",
// 			Email: "user1@example.com",
// 		},
// 		Active:  true,
// 		shop: 12,
// 	}

// 	return data, shopView
// }

// func testshopData1() (domain.shopData, view.shopDataView) {
// 	id := uuid.MustParse("00000000-0000-0000-0000-000000000001")

// 	data := domain.shopData{
// 		User: &domain.User{
// 			ID:    id,
// 			Name:  "User 1",
// 			Email: "user1@example.com",
// 		},
// 		Active:            true,
// 		shop:           12,
// 		Usedshop:       10,
// 		LastActiveRequest: nil,
// 	}

// 	shopView := view.shopDataView{
// 		ID: id.String(),
// 		User: &view.shopUserView{
// 			ID:    id,
// 			Name:  "User 1",
// 			Email: "user1@example.com",
// 		},
// 		Active:         true,
// 		shop:        12,
// 		Usedshop:    10,
// 		shopRequest: nil,
// 	}

// 	return data, shopView
// }

// func testshopData2() (domain.shopData, view.shopDataView) {
// 	id := uuid.MustParse("00000000-0000-0000-0000-000000000002")

// 	data := domain.shopData{
// 		User: &domain.User{
// 			ID:    id,
// 			Name:  "User 2",
// 			Email: "user2@example.com",
// 		},
// 		Active:      false,
// 		shop:     20,
// 		Usedshop: 13,
// 		LastActiveRequest: &domain.shopRequest{
// 			ID:      3,
// 			Status:  domain.shopRequestStatusPending,
// 			shop: 30,
// 			Remark:  "Some remark",
// 		},
// 	}

// 	shopView := view.shopDataView{
// 		ID: id.String(),
// 		User: &view.shopUserView{
// 			ID:    id,
// 			Name:  "User 2",
// 			Email: "user2@example.com",
// 		},
// 		Active:      false,
// 		shop:     20,
// 		Usedshop: 13,
// 		shopRequest: &view.shopRequestView{
// 			ID:      3,
// 			shop: 30,
// 			Status:  "Pending",
// 			Remark:  "Some remark",
// 		},
// 	}

// 	return data, shopView
// }
