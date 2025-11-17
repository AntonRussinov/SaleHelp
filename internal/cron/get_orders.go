package cron

import (
	"context"
	"fmt"
	"log/slog"
	"salehelp/internal/adapter/kaspi/portal/connector"
	"salehelp/internal/storage/pggorm/repository/merchant"

	"gorm.io/gorm"
)

// type GetOrdersRequest struct {
// 	PageNumber int
// 	PageSize   int
// 	Filter     struct {
// 		Orders struct {
// 			CreationDateGe    time.Time
// 			CreationDateLe    time.Time
// 			State             OrdersState
// 			Status            OrdersStatus
// 			DeliveryType      OrdersDeliveryType
// 			SignatureRequired bool
// 		}
// 	}
// 	Include struct {
// 		Orders string
// 	}
// }

func GetOrders(ctx context.Context, db *gorm.DB, log *slog.Logger) {
	log.Info("GetOrders started")

	// session := "1741341843.102.12897.790894|825e5f3659dba1ed7b5d7b2cbf5f1012"
	// sid := "cfeb3fe1-141f-4f02-bb37-72cf61adf722"
	merchRepo := merchant.NewRepository(log, db)
	api := connector.New(merchRepo)

	// orderRequest := models.GetOrdersRequest{
	// 	PageNumber: 0,
	// 	PageSize:   10,
	// 	Filter: models.OrderFilter{
	// 		Orders: models.Order{
	// 			CreationDateGe: time.Now().AddDate(0, 0, -3),
	// 			CreationDateLe: time.Now(),
	// 			State:          models.OrdersStateKaspiDelivery,
	// 		},
	// 	},
	// }

	products, err := api.GetAllProducts(ctx, "30345212")

	//orders, err := api.GetOrders(ctx, orderRequest)
	if err != nil {
		fmt.Println("Error getting orders: ", err)
		return
	}

	//for _, order := range orders.  {
	fmt.Println("Data: ", products.Data)
	fmt.Println("Data: ", products.Total)
	//}

	for _, data := range products.Data {
		fmt.Println("Included: ", data.MasterTitle)
	}
}
