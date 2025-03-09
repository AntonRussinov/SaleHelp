package models

type (
	ProductResponse struct {
		Data struct {
			Id            string            `json:"id"`
			Type          string            `json:"type"`
			Attributes    ProductAttributes `json:"attributes"`
			Relationships struct {
				MerchantProduct Relationship `json:"merchantProduct"`
			} `json:"relationships"`
			Links struct {
				Self string `json:"self"`
			} `json:"links"`
		} `json:"data"`
		Included []interface{} `json:"included"`
	}

	ProductAttributes struct {
		Code         string `json:"code"`
		Name         string `json:"name"`
		Manufacturer string `json:"manufacturer"`
	}
)
