package salamoonder

import "errors"

// Common Types

var ErrTaskNotReady = errors.New("task not ready")

type TaskPayload interface {
	BalanceTask | KasadaTask | TwitchScraperTask | TwitchPublicIntegrityTask | TwitchLocalIntegrityTask
}

type SolutionPayload interface {
	KasadaSolution | TwitchScraperSolution | TwitchPublicIntegritySolution | TwitchLocalIntegritySolution
}

type CreateTaskResponse struct {
	ErrorCode        int    `json:"error_code"`
	ErrorDescription string `json:"error_description"`
	TaskId           string `json:"taskId"`
}

type CreateTaskRequest[T TaskPayload] struct {
	APIKey string `json:"api_key"`
	Task   T      `json:"task"`
}

type TaskRequest struct {
	APIKey string `json:"api_key"`
	TaskId string `json:"taskId"`
}

type TaskResponse[S SolutionPayload] struct {
	ErrorId  int    `json:"error_id"`
	Solution S      `json:"solution"`
	Status   string `json:"status"`
}

// Balance

type BalanceTask struct {
	Type string `json:"type"` // "getBalance"
}

type BalanceResponse struct {
	ErrorCode        int    `json:"error_code"`
	ErrorDescription string `json:"error_description"`
	Wallet           string `json:"wallet"`
}

// Kasada Solver

type KasadaTask struct {
	Type   string `json:"type"` // "KasadaCaptchaSolver"
	PJS    string `json:"pjs"`
	CDOnly string `json:"cdOnly"` // "true/false"
}

type KasadaSolution struct {
	UserAgent string `json:"user-agent"`
	XCd       string `json:"x-kpsdk-cd"`
	XCr       string `json:"x-kpsdk-cr"`
	XCt       string `json:"x-kpsdk-ct"`
	XR        string `json:"x-kpsdk-r"`
	XV        string `json:"x-kpsdk-v"`
	XSt       string `json:"x-kpsdk-st"`
}

// Twitch Scraper

type TwitchScraperTask struct {
	Type string `json:"type"` // "TwitchScraper"
}

type TwitchScraperSolution struct {
	Biography      string `json:"biography"`
	ProfilePicture string `json:"profile_picture"`
	Username       string `json:"username"`
}

// Twitch Public Integrity

type TwitchPublicIntegrityTask struct {
	Type        string `json:"type"` // "Twitch_PublicIntegrity"
	AccessToken string `json:"access_token"`
	Proxy       string `json:"proxy"`
	DeviceID    string `json:"deviceId"` // optional
	ClientID    string `json:"clientId"` // optional
}

type TwitchPublicIntegritySolution struct {
	DeviceID       string `json:"device_id"`
	IntegrityToken string `json:"integrity_token"`
	Proxy          string `json:"proxy"`
	UserAgent      string `json:"user-agent"`
	ClientID       string `json:"client-id"`
}

// Twitch Local Integrity

type TwitchLocalIntegrityTask struct {
	Type     string `json:"type"` // "Twitch_LocalIntegrity"
	Proxy    string `json:"proxy"`
	DeviceID string `json:"deviceId"` // optional
	ClientID string `json:"clientId"` // optional
}

type TwitchLocalIntegritySolution struct {
	DeviceID       string `json:"device_id"`
	IntegrityToken string `json:"integrity_token"`
	Proxy          string `json:"proxy"`
	UserAgent      string `json:"user-agent"`
	ClientID       string `json:"client-id"`
}
