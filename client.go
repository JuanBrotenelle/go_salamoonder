package salamoonder

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
)

type client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

type Salamoonder interface {
	// Balance returns current wallet balance as a number.
	Balance() (float64, error)
	// KasadaCreate submits a Kasada task and returns taskId.
	KasadaCreate(pjs string, cdOnly bool) (string, error)
	// Kasada fetches a Kasada task result.
	Kasada(taskId string) (KasadaSolution, error)
	// TwitchScraperCreate submits a Twitch scraper task and returns taskId.
	TwitchScraperCreate() (string, error)
	// TwitchScraper fetches a Twitch scraper task result.
	TwitchScraper(taskId string) (TwitchScraperSolution, error)
	// TwitchPICreate submits a Twitch Public Integrity task and returns taskId. optional: deviceId, clientId
	TwitchPICreate(accessToken, proxy string, optional ...string) (string, error)
	// TwitchPI fetches a Twitch Public Integrity task result.
	TwitchPI(taskId string) (TwitchPublicIntegritySolution, error)
	// TwitchLICreate submits a Twitch Local Integrity task and returns taskId. optional: deviceId, clientId
	TwitchLICreate(proxy string, optional ...string) (string, error)
	// TwitchLI fetches a Twitch Local Integrity task result.
	TwitchLI(taskId string) (TwitchLocalIntegritySolution, error)
	// ContextBalance returns current wallet balance as a number. With context
	ContextBalance(ctx context.Context) (float64, error)
	// ContextKasadaCreate submits a Kasada task and returns taskId. With context
	ContextKasadaCreate(ctx context.Context, pjs string, cdOnly bool) (string, error)
	// ContextKasada fetches a Kasada task result. With context
	ContextKasada(ctx context.Context, taskId string) (KasadaSolution, error)
	// ContextTwitchScraperCreate submits a Twitch scraper task and returns taskId. With context
	ContextTwitchScraperCreate(ctx context.Context) (string, error)
	// ContextTwitchScraper fetches a Twitch scraper task result. With context
	ContextTwitchScraper(ctx context.Context, taskId string) (TwitchScraperSolution, error)
	// ContextTwitchPICreate submits a Twitch Public Integrity task and returns taskId. optional: deviceId, clientId. With context
	ContextTwitchPICreate(ctx context.Context, accessToken, proxy string, optional ...string) (string, error)
	// ContextTwitchPI fetches a Twitch Public Integrity task result. With context
	ContextTwitchPI(ctx context.Context, taskId string) (TwitchPublicIntegritySolution, error)
	// ContextTwitchLICreate submits a Twitch Local Integrity task and returns taskId. optional: deviceId, clientId. With context
	ContextTwitchLICreate(ctx context.Context, proxy string, optional ...string) (string, error)
	// ContextTwitchLI fetches a Twitch Local Integrity task result. With context
	ContextTwitchLI(ctx context.Context, taskId string) (TwitchLocalIntegritySolution, error)
}

func New(apiKey string) Salamoonder {
	return &client{
		baseURL:    "https://salamoonder.com/api",
		apiKey:     apiKey,
		httpClient: &http.Client{},
	}
}

// ContextsetBaseURL is an unexported helper used in tests to point the client to a test server.
func (c *client) setBaseURL(url string) { c.baseURL = url }

// ContextBalance returns current wallet balance as a number.
func (c *client) ContextBalance(ctx context.Context) (float64, error) {
	resp, err := createTask[BalanceTask, BalanceResponse](context.Background(), c, BalanceTask{Type: "getBalance"})
	if err != nil {
		return 0, err
	}
	if resp.ErrorCode != 0 {
		return 0, fmt.Errorf("api error %d: %s", resp.ErrorCode, resp.ErrorDescription)
	}
	amount, err := strconv.ParseFloat(resp.Wallet, 64)
	if err != nil {
		return 0, fmt.Errorf("parse wallet '%s': %w", resp.Wallet, err)
	}
	return amount, nil
}

func (c *client) Balance() (float64, error) {
	return c.ContextBalance(context.Background())
}

// ContextKasadaCreate submits a Kasada task and returns taskId.
func (c *client) ContextKasadaCreate(ctx context.Context, pjs string, cdOnly bool) (string, error) {
	resp, err := createTask[KasadaTask, CreateTaskResponse](context.Background(), c, KasadaTask{
		Type:   "KasadaCaptchaSolver",
		PJS:    pjs,
		CDOnly: boolToString(cdOnly),
	})
	if err != nil {
		return "", err
	}
	if resp.ErrorCode != 0 {
		return "", fmt.Errorf("api error %d: %s", resp.ErrorCode, resp.ErrorDescription)
	}
	return resp.TaskId, nil
}

// ContextKasada fetches a Kasada task result.
func (c *client) ContextKasada(ctx context.Context, taskId string) (KasadaSolution, error) {
	resp, err := taskResult[TaskResponse[KasadaSolution]](context.Background(), c, taskId)
	if err != nil {
		return KasadaSolution{}, err
	}
	if resp.ErrorId != 0 {
		return KasadaSolution{}, fmt.Errorf("api error %d", resp.ErrorId)
	}
	if resp.Status != "ready" {
		return KasadaSolution{}, ErrTaskNotReady
	}
	return resp.Solution, nil
}

func (c *client) KasadaCreate(pjs string, cdOnly bool) (string, error) {
	return c.ContextKasadaCreate(context.Background(), pjs, cdOnly)
}

func (c *client) Kasada(taskId string) (KasadaSolution, error) {
	return c.ContextKasada(context.Background(), taskId)
}

// ContextTwitchScraperCreate submits a Twitch scraper task and returns taskId. With context
func (c *client) ContextTwitchScraperCreate(ctx context.Context) (string, error) {
	resp, err := createTask[TwitchScraperTask, CreateTaskResponse](context.Background(), c, TwitchScraperTask{Type: "TwitchScraper"})
	if err != nil {
		return "", err
	}
	if resp.ErrorCode != 0 {
		return "", fmt.Errorf("api error %d: %s", resp.ErrorCode, resp.ErrorDescription)
	}
	return resp.TaskId, nil
}

// ContextTwitchScraper fetches a Twitch scraper task result. With context
func (c *client) ContextTwitchScraper(ctx context.Context, taskId string) (TwitchScraperSolution, error) {
	resp, err := taskResult[TaskResponse[TwitchScraperSolution]](context.Background(), c, taskId)
	if err != nil {
		return TwitchScraperSolution{}, err
	}
	if resp.ErrorId != 0 {
		return TwitchScraperSolution{}, fmt.Errorf("api error %d", resp.ErrorId)
	}
	if resp.Status != "ready" {
		return TwitchScraperSolution{}, ErrTaskNotReady
	}
	return resp.Solution, nil
}

func (c *client) TwitchScraperCreate() (string, error) {
	return c.ContextTwitchScraperCreate(context.Background())
}

func (c *client) TwitchScraper(taskId string) (TwitchScraperSolution, error) {
	return c.ContextTwitchScraper(context.Background(), taskId)
}

// ContextTwitchPICreate submits a Twitch Public Integrity task and returns taskId.
func (c *client) ContextTwitchPICreate(ctx context.Context, accessToken, proxy string, optional ...string) (string, error) {
	var deviceID, clientID string
	if len(optional) > 0 {
		deviceID = optional[0]
	}
	if len(optional) > 1 {
		clientID = optional[1]
	}
	resp, err := createTask[TwitchPublicIntegrityTask, CreateTaskResponse](context.Background(), c, TwitchPublicIntegrityTask{
		Type:        "Twitch_PublicIntegrity",
		AccessToken: accessToken,
		Proxy:       proxy,
		DeviceID:    deviceID,
		ClientID:    clientID,
	})
	if err != nil {
		return "", err
	}
	if resp.ErrorCode != 0 {
		return "", fmt.Errorf("api error %d: %s", resp.ErrorCode, resp.ErrorDescription)
	}
	return resp.TaskId, nil
}

// ContextTwitchPI fetches a Twitch Public Integrity task result. With context
func (c *client) ContextTwitchPI(ctx context.Context, taskId string) (TwitchPublicIntegritySolution, error) {
	resp, err := taskResult[TaskResponse[TwitchPublicIntegritySolution]](context.Background(), c, taskId)
	if err != nil {
		return TwitchPublicIntegritySolution{}, err
	}
	if resp.ErrorId != 0 {
		return TwitchPublicIntegritySolution{}, fmt.Errorf("api error %d", resp.ErrorId)
	}
	if resp.Status != "ready" {
		return TwitchPublicIntegritySolution{}, ErrTaskNotReady
	}
	return resp.Solution, nil
}

func (c *client) TwitchPICreate(accessToken, proxy string, optional ...string) (string, error) {
	return c.ContextTwitchPICreate(context.Background(), accessToken, proxy, optional...)
}

func (c *client) TwitchPI(taskId string) (TwitchPublicIntegritySolution, error) {
	return c.ContextTwitchPI(context.Background(), taskId)
}

// ContextTwitchLICreate submits a Twitch Local Integrity task and returns taskId.
func (c *client) ContextTwitchLICreate(ctx context.Context, proxy string, optional ...string) (string, error) {
	var deviceID, clientID string
	if len(optional) > 0 {
		deviceID = optional[0]
	}
	if len(optional) > 1 {
		clientID = optional[1]
	}
	resp, err := createTask[TwitchLocalIntegrityTask, CreateTaskResponse](context.Background(), c, TwitchLocalIntegrityTask{
		Type:     "Twitch_LocalIntegrity",
		Proxy:    proxy,
		DeviceID: deviceID,
		ClientID: clientID,
	})
	if err != nil {
		return "", err
	}
	if resp.ErrorCode != 0 {
		return "", fmt.Errorf("api error %d: %s", resp.ErrorCode, resp.ErrorDescription)
	}
	return resp.TaskId, nil
}

// ContextTwitchLI fetches a Twitch Local Integrity task result. With context
func (c *client) ContextTwitchLI(ctx context.Context, taskId string) (TwitchLocalIntegritySolution, error) {
	resp, err := taskResult[TaskResponse[TwitchLocalIntegritySolution]](context.Background(), c, taskId)
	if err != nil {
		return TwitchLocalIntegritySolution{}, err
	}
	if resp.ErrorId != 0 {
		return TwitchLocalIntegritySolution{}, fmt.Errorf("api error %d", resp.ErrorId)
	}
	if resp.Status != "ready" {
		return TwitchLocalIntegritySolution{}, ErrTaskNotReady
	}
	return resp.Solution, nil
}

func (c *client) TwitchLICreate(proxy string, optional ...string) (string, error) {
	return c.ContextTwitchLICreate(context.Background(), proxy, optional...)
}

func (c *client) TwitchLI(taskId string) (TwitchLocalIntegritySolution, error) {
	return c.ContextTwitchLI(context.Background(), taskId)
}
