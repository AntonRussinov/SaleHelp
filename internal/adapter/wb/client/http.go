package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// errRetry429 — после 429 и успешного sleep продолжаем цикл retry, не возвращаем успех.
var errRetry429 = errors.New("retry after 429")

func (c *Client) doJSON(
	ctx context.Context,
	method, url, apiKey string,
	cat RateCategory,
	in any,
	out any,
) error {

	var bodyBytes []byte
	var err error

	// 1) Готовим body ОДИН РАЗ (важно для retry)
	if in != nil {
		bodyBytes, err = json.Marshal(in)
		if err != nil {
			return fmt.Errorf("marshal request: %w", err)
		}
	}

	const maxAttempts = 5

	for attempt := 1; attempt <= maxAttempts; attempt++ {

		// 2) Rate limit — перед КАЖДОЙ попыткой
		if c.limiters != nil && apiKey != "" {
			if err := c.limiters.Get(apiKey, cat).Wait(ctx); err != nil {
				return err
			}
		}

		// 3) Новый request на каждую попытку
		var bodyReader io.Reader
		if bodyBytes != nil {
			bodyReader = bytes.NewReader(bodyBytes)
		}

		req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
		if err != nil {
			return fmt.Errorf("new request: %w", err)
		}

		if bodyBytes != nil {
			req.Header.Set("Content-Type", "application/json")
		}
		if apiKey != "" {
			req.Header.Set("Authorization", apiKey)
		}

		// 4) Выполняем запрос
		resp, err := c.httpClient.Do(req)
		if err != nil {
			return fmt.Errorf("do request: %w", err)
		}

		// ВАЖНО: body закрываем в этой же итерации
		func() {
			defer resp.Body.Close()

			// 5) 429 — уважаем WB и пробуем ещё раз
			if resp.StatusCode == http.StatusTooManyRequests {
				if attempt == maxAttempts {
					err = fmt.Errorf("rate limited after %d attempts", maxAttempts)
					return
				}

				retryAfter := parseRetryAfter(resp.Header)
				if retryAfter <= 0 {
					retryAfter = 20 // safe fallback
				}
				if sleepErr := sleepWithContext(ctx, time.Duration(retryAfter)*time.Second); sleepErr != nil {
					err = sleepErr
					return
				}
				err = errRetry429 // не возвращать успех — продолжить цикл и повторить запрос
				return
			}

			// 6) Любой другой не-2xx — ошибка без retry
			if resp.StatusCode < 200 || resp.StatusCode >= 300 {
				b, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
				err = &HTTPError{
					StatusCode: resp.StatusCode,
					URL:        url, // переменная url из аргументов doJSON
					Body:       strings.TrimSpace(string(b)),
				}
				return
			}

			// 7) Успех
			if out == nil {
				err = nil
				return
			}

			err = json.NewDecoder(resp.Body).Decode(out)
			if err != nil {
				err = fmt.Errorf("decode response: %w", err)
			}
		}()

		if err == nil {
			return nil
		}
		if errors.Is(err, errRetry429) {
			continue // 429 + sleep успешен — повторяем запрос
		}
		return err // 4xx/5xx или rate limited на последней попытке
	}

	return fmt.Errorf("request failed after %d attempts", maxAttempts)
}

func parseRetryAfter(h http.Header) int {
	if v := h.Get("X-Ratelimit-Retry"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	if v := h.Get("Retry-After"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return 0
}

func sleepWithContext(ctx context.Context, d time.Duration) error {
	t := time.NewTimer(d)
	defer t.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-t.C:
		return nil
	}
}
