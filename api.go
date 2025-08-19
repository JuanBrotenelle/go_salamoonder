package salamoonder

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
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

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status: %s", resp.Status)
	}

	if err := json.NewDecoder(resp.Body).Decode(responseDest); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}

	return nil
}

// createTask submits a task payload and decodes the response into the caller-provided generic type R.
// T is the concrete task payload type; R is the expected response type from the endpoint.
func createTask[T TaskPayload, R any](ctx context.Context, c *client, task T) (R, error) {
	var respBody R
	rawReq := CreateTaskRequest[T]{
		APIKey: c.apiKey,
		Task:   task,
	}
	if err := c.postJSON(ctx, "/createTask", rawReq, &respBody); err != nil {
		return respBody, fmt.Errorf("createTask: %w", err)
	}
	return respBody, nil
}

// taskResult requests a task result by id and decodes it into R.
func taskResult[R any](ctx context.Context, c *client, taskId string) (R, error) {
	var respBody R
	rawReq := TaskRequest{
		APIKey: c.apiKey,
		TaskId: taskId,
	}
	if err := c.postJSON(ctx, "/getTaskResult", rawReq, &respBody); err != nil {
		return respBody, fmt.Errorf("taskResult: %w", err)
	}
	return respBody, nil
}
