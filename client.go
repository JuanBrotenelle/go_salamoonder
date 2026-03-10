package salamoonder

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

type Client struct {
	*client
}

func New(apiKey string, httpClient *http.Client) (*Client, error) {
	if apiKey == "" {
		return nil, ErrNoApiKey
	}

	if httpClient == nil {
		httpClient = &http.Client{}
	}

	return &Client{
		client: &client{
			baseURL:    "https://salamoonder.com/api",
			apiKey:     apiKey,
			httpClient: httpClient,
		},
	}, nil
}

func (c *Client) Balance(ctx context.Context) (*CreateTaskBalanceResult, error) {
	req := CreateTaskRequest{
		ApiKey: c.apiKey,
		Task: map[string]interface{}{
			"type": "getBalance",
		},
	}

	var result CreateTaskBalanceResult
	if err := c.postJSON(ctx, "/createTask", req, &result); err != nil {
		return nil, err
	}

	if result.ErrorCode != 0 {
		return &result, &APIError{
			StatusCode: http.StatusOK,
			Msg:        result.ErrorDescription,
		}
	}

	return &result, nil
}

func (c *Client) CreateTask(ctx context.Context, options any) (*CreateTaskResult, error) {
	switch opts := options.(type) {
	case KasadaStandardOptions:
		return createTaskGeneric(c, ctx, opts)
	case KasadaPayloadOptions:
		return createTaskGeneric(c, ctx, opts)
	case AkamaiWebOptions:
		return createTaskGeneric(c, ctx, opts)
	case AkamaiSBSDOptions:
		return createTaskGeneric(c, ctx, opts)
	case Reese84Options:
		return createTaskGeneric(c, ctx, opts)
	case UutmvcOptions:
		return createTaskGeneric(c, ctx, opts)
	case DataDomeInterstitialOptions:
		return createTaskGeneric(c, ctx, opts)
	case DataDomeSliderOptions:
		return createTaskGeneric(c, ctx, opts)
	case TwitchScraperOptions:
		return createTaskGeneric(c, ctx, opts)
	case TwitchIntegrityOptions:
		return createTaskGeneric(c, ctx, opts)
	default:
		return nil, &MethodError{
			OptionsValue: options,
		}
	}
}

func (c *Client) Task(ctx context.Context, taskId string) (*TaskResultRaw, error) {
	var result TaskResultRaw
	req := TaskRequest{
		APIKey: c.apiKey,
		TaskId: taskId,
	}
	if err := c.postJSON(ctx, "/getTaskResult", req, &result); err != nil {
		return nil, err
	}
	if result.ErrorId != 0 {
		return &result, &APIError{
			StatusCode: http.StatusOK,
			TaskId:     taskId,
			Msg:        result.Status,
		}
	}
	return &result, nil
}

func createTaskGeneric[TO TaskOptions](c *Client, ctx context.Context, options TO) (*CreateTaskResult, error) {
	taskType := getTaskTypeFromOptions(options)

	optionsJSON, err := json.Marshal(options)
	if err != nil {
		return nil, fmt.Errorf("marshal options: %w", err)
	}

	var optionsMap map[string]interface{}
	if err := json.Unmarshal(optionsJSON, &optionsMap); err != nil {
		return nil, fmt.Errorf("unmarshal options: %w", err)
	}

	taskPayload := map[string]interface{}{
		"type": taskType,
	}

	for k, v := range optionsMap {
		taskPayload[k] = v
	}

	req := CreateTaskRequest{
		ApiKey: c.apiKey,
		Task:   taskPayload,
	}

	var result CreateTaskResult
	if err := c.postJSON(ctx, "/createTask", req, &result); err != nil {
		return nil, err
	}

	if result.ErrorCode != 0 {
		return &result, &APIError{
			StatusCode: http.StatusOK,
			TaskId:     result.TaskId,
			Msg:        result.ErrorDescription,
		}
	}

	return &result, nil
}

func GetTaskResult[TS TaskSolution](c *Client, ctx context.Context, taskId string) (*TaskResult[TS], error) {
	var result TaskResult[TS]
	req := TaskRequest{
		APIKey: c.apiKey,
		TaskId: taskId,
	}
	if err := c.postJSON(ctx, "/getTaskResult", req, &result); err != nil {
		return nil, err
	}
	if result.ErrorId != 0 {
		return &result, &APIError{
			StatusCode: http.StatusOK,
			TaskId:     taskId,
			Msg:        result.Status,
		}
	}
	return &result, nil
}

func (c *client) setBaseURL(url string) {
	c.baseURL = url
}

func getTaskTypeFromOptions(opts any) string {
	switch any(opts).(type) {
	case KasadaStandardOptions:
		return "KasadaCaptchaSolver"
	case KasadaPayloadOptions:
		return "KasadaPayloadSolver"
	case AkamaiWebOptions:
		return "AkamaiWebSolver"
	case AkamaiSBSDOptions:
		return "AkamaiSBSDSolver"
	case Reese84Options:
		return "IncapsulaReese84Solver"
	case UutmvcOptions:
		return "IncapsulaUTMVCSolver"
	case DataDomeInterstitialOptions:
		return "DataDomeInterstitialSolver"
	case DataDomeSliderOptions:
		return "DataDomeSliderSolver"
	case TwitchScraperOptions:
		return "Twitch_Scraper"
	case TwitchIntegrityOptions:
		return "Twitch_PublicIntegrity"
	default:
		return "unknown"
	}
}
