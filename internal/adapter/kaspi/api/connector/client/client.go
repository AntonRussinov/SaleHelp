package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type (
	Client struct {
		authToken string
	}

	method string
)

const (
	GET  method = "GET"
	POST method = "POST"
)

func New() *Client {
	return &Client{}
}

func (a *Client) Do(ctx context.Context, method method, path string, params url.Values, authToken string, result interface{}) error {
	req, err := http.NewRequestWithContext(ctx, string(method), fmt.Sprintf("%s/%s", "https://kaspi.kz/shop/api/v2", path), nil)
	if err != nil {
		return err
	}

	if params != nil {
		req.URL.RawQuery = params.Encode()
	}

	req.Header.Set("Accept", "*/*")
	req.Header.Set("Content-Type", "application/vnd.api+json")
	req.Header.Set("X-Auth-Token", authToken)

	httpClient := &http.Client{}

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if result != nil {
		return json.NewDecoder(resp.Body).Decode(result)
	}
	if _, err = io.Copy(io.Discard, resp.Body); err != nil {
		return err
	}

	return nil
}
