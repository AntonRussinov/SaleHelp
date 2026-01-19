package client

import "fmt"

type HTTPError struct {
	StatusCode int
	URL        string
	Body       string
}

func (e *HTTPError) Error() string {
	if e.Body != "" {
		return fmt.Sprintf("wb http error: status=%d url=%s body=%s", e.StatusCode, e.URL, e.Body)
	}
	return fmt.Sprintf("wb http error: status=%d url=%s", e.StatusCode, e.URL)
}
