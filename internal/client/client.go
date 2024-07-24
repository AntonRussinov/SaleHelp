package client

import (
	"github.com/go-resty/resty/v2"
)

var restyClient = resty.New()

func GetClientData(clientID string) (*resty.Response, error) {
	resp, err := restyClient.R().
		SetPathParams(map[string]string{"clientID": clientID}).
		Get("https://api.external.com/clients/{clientID}")
	return resp, err
}
