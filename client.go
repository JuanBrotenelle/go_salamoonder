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

// Salamoonder describes the client contract for interacting with the salamoonder.com API.
// Each method corresponds to creating a task via createTask and retrieving its result via getTaskResult.
// Methods with the Context prefix accept context.Context and allow managing timeouts and request cancellation.
type Salamoonder interface {
	// Balance returns the current balance of the user's wallet.
	//
	// Uses the `getBalance` task.
	// Cost: free.
	//
	// Returns:
	//   - float64 — balance in the service's currency.
	//   - error   — if the request fails or the server returns an error.
	Balance() (float64, error)

	// KasadaCreate creates a task to solve Kasada Captcha (type: "KasadaCaptchaSolver").
	//
	// Parameters:
	//   - pjs    — URL to the p.js script.
	//   - cdOnly — if true, only CD tokens are returned.
	//
	// Cost: 0.002 cents per request.
	//
	// Returns:
	//   - taskId — unique task identifier.
	//   - error  — if task creation fails.
	KasadaCreate(pjs string, cdOnly bool) (string, error)

	// Kasada retrieves the result of a previously created Kasada task.
	//
	// Parameters:
	//   - taskId — task identifier issued by KasadaCreate.
	//
	// Returns:
	//   - KasadaSolution — object containing headers of the form x-kpsdk-*.
	//   - error          — if the task is not ready or an error occurs.
	Kasada(taskId string) (KasadaSolution, error)

	// TwitchScraperCreate creates a Twitch scraper task (type: "Twitch_Scraper").
	//
	// Cost: 0.0001 cents per request.
	//
	// Returns:
	//   - taskId — unique task identifier.
	//   - error  — if an error occurs.
	TwitchScraperCreate() (string, error)

	// TwitchScraper retrieves the result of a Twitch scraper task.
	//
	// Parameters:
	//   - taskId — task identifier.
	//
	// Returns:
	//   - TwitchScraperSolution — contains username, biography, profile_picture.
	//   - error                 — if the task is not ready or an error occurs.
	TwitchScraper(taskId string) (TwitchScraperSolution, error)

	// PublicIntegrityCreate creates a Twitch Public Integrity task (type: "Twitch_PublicIntegrity").
	//
	// Requires own proxies (IP-auth proxies are not supported).
	//
	// Parameters:
	//   - accessToken — Twitch OAuth token.
	//   - proxy       — string in the format "user:pass@ip:port".
	//   - optional    — deviceId, clientId (optional).
	//
	// Cost: 0.002 cents per request.
	//
	// Returns:
	//   - taskId — unique task identifier.
	//   - error  — if task creation fails.
	PublicIntegrityCreate(accessToken, proxy string, optional ...string) (string, error)

	// PublicIntegrity retrieves the result of a Twitch Public Integrity task.
	//
	// Parameters:
	//   - taskId — task identifier.
	//
	// Returns:
	//   - TwitchPublicIntegritySolution — integrity_token, device_id, client-id, user-agent.
	//   - error                         — if the task is not ready or an error occurs.
	PublicIntegrity(taskId string) (TwitchPublicIntegritySolution, error)

	// LocalIntegrityCreate creates a Twitch Local Integrity task (type: "Twitch_LocalIntegrity").
	//
	// Requires own proxies.
	//
	// Parameters:
	//   - proxy    — string in the format "user:pass@ip:port".
	//   - optional — deviceId, clientId (optional).
	//
	// Cost: 0.002 cents per request.
	//
	// Returns:
	//   - taskId — unique task identifier.
	//   - error  — if task creation fails.
	LocalIntegrityCreate(proxy string, optional ...string) (string, error)

	// LocalIntegrity retrieves the result of a Twitch Local Integrity task.
	//
	// Parameters:
	//   - taskId — task identifier.
	//
	// Returns:
	//   - TwitchLocalIntegritySolution — integrity_token, device_id, client-id, user-agent.
	//   - error                        — if the task is not ready or an error occurs.
	LocalIntegrity(taskId string) (TwitchLocalIntegritySolution, error)

	// --- Methods with context.Context ---

	// ContextBalance is similar to Balance but with context.
	ContextBalance(ctx context.Context) (float64, error)

	// ContextKasadaCreate is similar to KasadaCreate but with context.
	ContextKasadaCreate(ctx context.Context, pjs string, cdOnly bool) (string, error)

	// ContextKasada is similar to Kasada but with context.
	ContextKasada(ctx context.Context, taskId string) (KasadaSolution, error)

	// ContextTwitchScraperCreate is similar to TwitchScraperCreate but with context.
	ContextTwitchScraperCreate(ctx context.Context) (string, error)

	// ContextTwitchScraper is similar to TwitchScraper but with context.
	ContextTwitchScraper(ctx context.Context, taskId string) (TwitchScraperSolution, error)

	// ContextPublicIntegrityCreate is similar to PublicIntegrityCreate but with context.
	ContextPublicIntegrityCreate(ctx context.Context, accessToken, proxy string, optional ...string) (string, error)

	// ContextPublicIntegrity is similar to PublicIntegrity but with context.
	ContextPublicIntegrity(ctx context.Context, taskId string) (TwitchPublicIntegritySolution, error)

	// ContextLocalIntegrityCreate is similar to LocalIntegrityCreate but with context.
	ContextLocalIntegrityCreate(ctx context.Context, proxy string, optional ...string) (string, error)

	// ContextLocalIntegrity is similar to LocalIntegrity but with context.
	ContextLocalIntegrity(ctx context.Context, taskId string) (TwitchLocalIntegritySolution, error)
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

// ContextPublicIntegrityCreate submits a Twitch Public Integrity task and returns taskId.
func (c *client) ContextPublicIntegrityCreate(ctx context.Context, accessToken, proxy string, optional ...string) (string, error) {
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

// ContextPublicIntegrity fetches a Twitch Public Integrity task result. With context
func (c *client) ContextPublicIntegrity(ctx context.Context, taskId string) (TwitchPublicIntegritySolution, error) {
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

func (c *client) PublicIntegrityCreate(accessToken, proxy string, optional ...string) (string, error) {
	return c.ContextPublicIntegrityCreate(context.Background(), accessToken, proxy, optional...)
}

func (c *client) PublicIntegrity(taskId string) (TwitchPublicIntegritySolution, error) {
	return c.ContextPublicIntegrity(context.Background(), taskId)
}

// ContextLocalIntegrityCreate submits a Twitch Local Integrity task and returns taskId.
func (c *client) ContextLocalIntegrityCreate(ctx context.Context, proxy string, optional ...string) (string, error) {
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

// ContextLocalIntegrity fetches a Twitch Local Integrity task result. With context
func (c *client) ContextLocalIntegrity(ctx context.Context, taskId string) (TwitchLocalIntegritySolution, error) {
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

func (c *client) LocalIntegrityCreate(proxy string, optional ...string) (string, error) {
	return c.ContextLocalIntegrityCreate(context.Background(), proxy, optional...)
}

func (c *client) LocalIntegrity(taskId string) (TwitchLocalIntegritySolution, error) {
	return c.ContextLocalIntegrity(context.Background(), taskId)
}
