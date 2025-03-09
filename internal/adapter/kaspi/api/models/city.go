package models

type (
	CitiesResponse struct {
		Data []struct {
			Id            string         `json:"id"`
			Type          string         `json:"type"`
			Attributes    CityAttributes `json:"attributes"`
			Relationships struct {
			} `json:"relationships"`
			Links struct {
				Self string `json:"self"`
			} `json:"links"`
		} `json:"data"`
		Included []interface{} `json:"included"`
	}

	CityResponse struct {
		Data struct {
			Id            string         `json:"id"`
			Type          string         `json:"type"`
			Attributes    CityAttributes `json:"attributes"`
			Relationships struct {
			} `json:"relationships"`
			Links struct {
				Self string `json:"self"`
			} `json:"links"`
		} `json:"data"`
		Included []interface{} `json:"included"`
	}

	CityAttributes struct {
		Code   string `json:"code"`
		Name   string `json:"name"`
		Active bool   `json:"active"`
	}
)
