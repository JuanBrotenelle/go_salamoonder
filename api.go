package salamoonder

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (c *client) postJSON(ctx context.Context, path string, requestBody any, responseDest any) error {
	url := fmt.Sprintf("%s%s", c.baseURL, path)

	payload, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("http do: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode == http.StatusBadRequest {
		var apiErr struct {
			ErrorCode        int    `json:"error_code"`
			ErrorDescription string `json:"error_description"`
		}
		if jsonErr := json.Unmarshal(body, &apiErr); jsonErr == nil && apiErr.ErrorDescription != "" {
			return &APIError{
				StatusCode: http.StatusBadRequest,
				Msg:        apiErr.ErrorDescription,
			}
		}
		return &APIError{
			StatusCode: http.StatusBadRequest,
			Msg:        string(body),
		}
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}

	if err := json.Unmarshal(body, responseDest); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}

	return nil
}
