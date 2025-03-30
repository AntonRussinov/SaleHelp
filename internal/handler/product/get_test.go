package product_test

// import (
// 	"encoding/json"
// 	"net/http"
// 	"salehelp/domain"
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/mock"
// 	"github.com/stretchr/testify/require"
// )

// func TestGetUSer(t *testing.T) {
// 	t.Parallel()

// 	data1, view1 := testProductData1()

// 	response, err := json.Marshal(view1)
// 	require.NoError(t, err)

// 	testcases := []struct {
// 		name          string
// 		paramValue    string
// 		serviceCalled bool
// 		serviceReturn *domain.ProductData
// 		serviceError  error
// 		wantErr       bool
// 	}{
// 		{
// 			name:          "success",
// 			paramValue:    data1.User.ID.String(),
// 			serviceCalled: true,
// 			serviceReturn: &data1,
// 			wantErr:       false,
// 		},
// 		{
// 			name:          "wrong param",
// 			paramValue:    "1",
// 			serviceCalled: false,
// 			wantErr:       true,
// 		},
// 		{
// 			name:          "service error",
// 			paramValue:    data1.User.ID.String(),
// 			serviceCalled: true,
// 			serviceReturn: nil,
// 			serviceError:  assert.AnError,
// 			wantErr:       true,
// 		},
// 	}

// 	for _, tt := range testcases {
// 		t.Run(tt.name, func(t *testing.T) {
// 			t.Parallel()

// 			test := setupHandlerTest(t, "GET", nil)

// 			if tt.serviceCalled {
// 				test.service.EXPECT().GetUserProductData(mock.Anything, data1.User.ID).
// 					Return(tt.serviceReturn, tt.serviceError).Once()
// 			}

// 			test.SetParam("id", tt.paramValue)
// 			err := test.handler.GetUserProductByUserID(test.c)

// 			if tt.wantErr {
// 				require.Error(t, err)
// 			} else {
// 				require.NoError(t, err)
// 				require.Equal(t, http.StatusOK, test.rec.Code)
// 				require.JSONEqf(t, string(response), test.rec.Body.String(), "error message %s", "formatted")
// 			}
// 		})
// 	}
// }
