package models

import (
	"errors"
	"net/url"
	"strconv"
	"time"
)

type (
	OrderEntriesResponse struct {
		Data []struct {
			ID            string               `json:"id"`
			Type          string               `json:"type"`
			Attributes    OrderEntryAttributes `json:"attributes"`
			Relationships struct {
				Product                Relationship `json:"product"`
				DeliveryPointOfService Relationship `json:"deliveryPointOfService"`
			} `json:"relationships"`
			Links struct {
				Self string `json:"self"`
			} `json:"links"`
		} `json:"data"`
		Included []interface{} `json:"included"`
	}

	OrderEntryResponse struct {
		Data struct {
			ID            string               `json:"id"`
			Type          string               `json:"type"`
			Attributes    OrderEntryAttributes `json:"attributes"`
			Relationships struct {
				Product                Relationship `json:"product"`
				DeliveryPointOfService Relationship `json:"deliveryPointOfService"`
			} `json:"relationships"`
			Links struct {
				Self string `json:"self"`
			} `json:"links"`
		} `json:"data"`
		Included []interface{} `json:"included"`
	}

	OrderEntryAttributes struct {
		EntryNumber  int     `json:"entryNumber"`
		Quantity     int     `json:"quantity"`
		TotalPrice   float64 `json:"totalPrice"`
		DeliveryCost float64 `json:"deliveryCost"`
		BasePrice    float64 `json:"basePrice"`
	}

	Relationship struct {
		Links Links `json:"links"`
		Data  Data  `json:"data"`
	}

	EntriesRelationship struct { //check it to delete it and just replace with Relationship
		Links Links `json:"links"`
	}

	Data struct {
		ID   string `json:"id"`
		Type string `json:"type"`
	}

	Links struct {
		Self    string `json:"self"`
		Related string `json:"related"`
	}

	GetOrdersRequest struct {
		PageNumber int
		PageSize   int
		Filter     OrderFilter
		Include    struct {
			Orders string
		}
	}

	OrderFilter struct {
		Orders Order
	}

	Order struct {
		CreationDateGe    time.Time
		CreationDateLe    time.Time
		State             OrdersState
		Status            OrdersStatus
		DeliveryType      OrdersDeliveryType
		SignatureRequired bool
	}

	OrdersResponse struct {
		Data []struct {
			ID            string          `json:"id"`
			Type          string          `json:"type"`
			Attributes    OrderAttributes `json:"attributes"`
			Relationships struct {
				Entries EntriesRelationship `json:"entries"`
				User    Relationship        `json:"user"`
			} `json:"relationships"`
			Links struct {
				Self string `json:"self"`
			} `json:"links"`
		} `json:"data"`
		Included []struct {
			ID         string `json:"id"`
			Type       string `json:"type"`
			Attributes struct {
				FirstName string `json:"firstName"`
				LastName  string `json:"lastName"`
				CellPhone string `json:"cellPhone"`
			} `json:"attributes"`

			Relationships struct {
			} `json:"relationships"`

			Links struct {
				Self string `json:"self"`
			} `json:"links"`
		} `json:"included"`
		Meta struct {
			PageCount  int `json:"pageCount"`
			TotalCount int `json:"totalCount"`
		} `json:"meta"`
	}

	OrderByCodeResponse struct {
		Data []struct {
			ID            string          `json:"id"`
			Type          string          `json:"type"`
			Attributes    OrderAttributes `json:"attributes"`
			Relationships struct {
				Entries EntriesRelationship `json:"entries"`
				User    Relationship        `json:"user"`
			} `json:"relationships"`
			Links struct {
				Self string `json:"self"`
			} `json:"links"`
		} `json:"data"`
		Included []interface{} `json:"included"`
		Meta     struct {
			PageCount  int `json:"pageCount"`
			TotalCount int `json:"totalCount"`
		} `json:"meta"`
	}

	OrderAttributes struct {
		State      OrdersState  `json:"state"`
		Status     OrdersStatus `json:"status"`
		Code       string       `json:"code"`
		TotalPrice float64      `json:"totalPrice"`
		Customer   struct {
			ID        string `json:"id"`
			Name      string `json:"name"`
			FirstName string `json:"firstName"`
			LastName  string `json:"lastName"`
			CellPhone string `json:"cellPhone"`
		} `json:"customer"`

		SignatureRequired   bool `json:"signatureRequired"`
		PreOrder            bool `json:"preOrder"`
		Express             bool `json:"express"`
		ReturnedToWarehouse bool `json:"returnedToWarehouse"`

		DeliveryMode    OrderDeliveryMode `json:"deliveryMode"`
		IsKaspiDelivery bool              `json:"isKaspiDelivery"`
		Latitude        float64           `json:"latitude"`
		Longitude       float64           `json:"longitude"`
		Apartment       string            `json:"apartment"`
		DeliveryAddress struct {
			StreetName       string      `json:"streetName"`
			StreetNumber     string      `json:"streetNumber"`
			Town             string      `json:"town"`
			District         interface{} `json:"district"`
			Building         interface{} `json:"building"`
			FormattedAddress string      `json:"formattedAddress"`
			Latitude         float64     `json:"latitude"`
			Longitude        float64     `json:"longitude"`
		} `json:"deliveryAddress"`
		KaspiDelivery struct {
			Waybill                         string `json:"waybill"`
			CourierTransmissionDate         int64  `json:"courierTransmissionDate"`
			CourierTransmissionPlanningDate int64  `json:"courierTransmissionPlanningDate"`
		} `json:"kaspiDelivery"`

		CreditTerm int `json:"creditTerm"`

		PaymentMode OrderPaymentMode `json:"paymentMode"`

		CreationDate                    int64 `json:"creationDate"`
		ApprovedByBankDate              int64 `json:"approvedByBankDate"`
		PlannedDeliveryDate             int64 `json:"plannedDeliveryDate"`
		ReservationDate                 int64 `json:"reservationDate"`
		CourierTransmissionPlanningDate int64 `json:"courierTransmissionPlanningDate"`
		CourierTransmissionDate         int64 `json:"courierTransmissionDate"`

		DeliveryCost          float64 `json:"deliveryCost"`
		DeliveryCostForSeller float64 `json:"deliveryCostForSeller"`

		Waybill       string `json:"waybill"`
		WaybillNumber string `json:"waybillNumber"`

		Category string `json:"category"`
	}

	OrderDeliveryMode  string
	OrdersDeliveryType string
	OrderPaymentMode   string
	OrdersState        string
	OrdersStatus       string
)

const (
	OrderDeliveryModeLocal          OrderDeliveryMode = "DELIVERY_LOCAL"
	OrderDeliveryModePickup         OrderDeliveryMode = "DELIVERY_PICKUP"
	OrderDeliveryModeRegionalPickup OrderDeliveryMode = "DELIVERY_REGIONAL_PICKUP"
	OrderDeliveryModeRegionalTodoor OrderDeliveryMode = "DELIVERY_REGIONAL_TODOOR"

	OrdersDeliveryTypePickup   OrdersDeliveryType = "PICKUP"
	OrdersDeliveryTypeDelivery OrdersDeliveryType = "DELIVERY"

	OrderPaymentModePayWithCredit OrderPaymentMode = "PAY_WITH_CREDIT"
	OrderPaymentModePrepaid       OrderPaymentMode = "PREPAID"

	OrdersStateNew           OrdersState = "NEW"
	OrdersStateSignRequired  OrdersState = "SIGN_REQUIRED"
	OrdersStatePickup        OrdersState = "PICKUP"
	OrdersStateDelivery      OrdersState = "DELIVERY"
	OrdersStateKaspiDelivery OrdersState = "KASPI_DELIVERY"
	OrdersStateArchive       OrdersState = "ARCHIVE"

	OrdersStatusApprovedByBank               OrdersStatus = "APPROVED_BY_BANK"
	OrdersStatusAcceptedByMerchant           OrdersStatus = "ACCEPTED_BY_MERCHANT"
	OrdersStatusCompleted                    OrdersStatus = "COMPLETED"
	OrdersStatusCancelled                    OrdersStatus = "CANCELLED"
	OrdersStatusCancelling                   OrdersStatus = "CANCELLING"
	OrdersStatusKaspiDeliveryReturnRequested OrdersStatus = "KASPI_DELIVERY_RETURN_REQUESTED"
	OrdersStatusReturnAcceptedByMerchant     OrdersStatus = "RETURN_ACCEPTED_BY_MERCHANT"
	OrdersStatusReturned                     OrdersStatus = "RETURNED"
)

var (
	ErrPageNumber                                    = errors.New("page number must be greater than 0")
	ErrPageSize                                      = errors.New("page size must be greater than 0")
	ErrPageSizeLimit                                 = errors.New("page size limit is 100")
	ErrFilterOrdersState                             = errors.New("filter orders state must be set")
	ErrFilterOrdersCreationDateGe                    = errors.New("filter orders creation date ge must be greater than 0")
	ErrFilterOrdersStatePickupDeliveryTypePickup     = errors.New("filter orders state pickup and delivery type pickup are incompatible")
	ErrFilterOrdersStateDeliveryDeliveryTypeDelivery = errors.New("filter orders state delivery and delivery type delivery are incompatible")
	ErrFilterOrdersStateSignRequired                 = errors.New("filter orders state sign required cannot be used without filter orders state sign required")
)

func (p *GetOrdersRequest) ToUrlValues() (url.Values, error) {
	if p.PageNumber < 0 {
		return nil, ErrPageNumber
	}
	if p.PageSize <= 0 {
		return nil, ErrPageSize
	}
	if p.PageSize > 100 {
		return nil, ErrPageSizeLimit
	}
	if p.Filter.Orders.CreationDateGe.Unix() <= 0 {
		return nil, ErrFilterOrdersCreationDateGe
	}
	if p.Filter.Orders.State == "" {
		return nil, ErrFilterOrdersState
	}
	if p.Filter.Orders.State == OrdersStatePickup && p.Filter.Orders.DeliveryType == OrdersDeliveryTypePickup {
		return nil, ErrFilterOrdersStatePickupDeliveryTypePickup
	}
	if p.Filter.Orders.State == OrdersStateDelivery && p.Filter.Orders.DeliveryType == OrdersDeliveryTypeDelivery {
		return nil, ErrFilterOrdersStateDeliveryDeliveryTypeDelivery
	}
	if p.Filter.Orders.SignatureRequired && p.Filter.Orders.State != OrdersStateSignRequired {
		return nil, ErrFilterOrdersStateSignRequired
	}

	params := make(url.Values)
	params.Add("page[number]", strconv.Itoa(p.PageNumber))
	params.Add("page[size]", strconv.Itoa(p.PageSize))
	params.Add("filter[orders][creationDate][$ge]", strconv.FormatInt(p.Filter.Orders.CreationDateGe.UnixMilli(), 10))
	if p.Filter.Orders.CreationDateLe.Unix() > 0 {
		params.Add("filter[orders][creationDate][$le]", strconv.FormatInt(p.Filter.Orders.CreationDateLe.UnixMilli(), 10))
	}
	params.Add("filter[orders][state]", string(p.Filter.Orders.State))
	if p.Filter.Orders.Status != "" {
		params.Add("filter[orders][status]", string(p.Filter.Orders.Status))
	}
	if p.Filter.Orders.DeliveryType != "" {
		params.Add("filter[orders][deliveryType]", string(p.Filter.Orders.DeliveryType))
	}
	if p.Filter.Orders.SignatureRequired {
		params.Add("filter[orders][signatureRequired]", strconv.FormatBool(p.Filter.Orders.SignatureRequired))
	}
	//if p.Include.Orders != "" {
	//	params.Add("include[orders]", p.Include.Orders)
	//}
	return params, nil
}
