package models

type (
	PointOfServiceResponse struct {
		Data struct {
			Id            string                   `json:"id"`
			Type          string                   `json:"type"`
			Attributes    PointOfServiceAttributes `json:"attributes"`
			Relationships struct {
				City Relationship `json:"city"`
			} `json:"relationships"`
			Links struct {
				Self string `json:"self"`
			} `json:"links"`
		} `json:"data"`
		Included []interface{} `json:"included"`
	}

	PointOfServiceAttributes struct {
		DisplayName string `json:"displayName"`
		Address     struct {
			StreetName       string      `json:"streetName"`
			StreetNumber     string      `json:"streetNumber"`
			Town             string      `json:"town"`
			District         interface{} `json:"district"`
			Building         interface{} `json:"building"`
			FormattedAddress string      `json:"formattedAddress"`
		} `json:"address"`
	}
)
