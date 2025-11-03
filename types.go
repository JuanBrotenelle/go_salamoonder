package salamoonder

import "encoding/json"

type TaskOptions interface {
	KasadaOptions | Reese84Options | UutmvcOptions | TwitchScraperOptions | TwitchIntegrityOptions
}

type TaskSolution interface {
	KasadaSolution | Reese84Solution | Reese84SubmitPayloadSolution | UutmvcSolution | TwitchScraperSolution | TwitchIntegritySolution
}

type CreateTaskRequest struct {
	ApiKey string      `json:"api_key"`
	Task   interface{} `json:"task"`
}

type CreateTaskResult struct {
	ErrorCode        int    `json:"error_code"`
	ErrorDescription string `json:"error_description"`
	TaskId           string `json:"taskId"`
}

// https://apidocs.salamoonder.com/tasks/get-balance
type CreateTaskBalanceResult struct {
	ErrorCode        int    `json:"error_code"`
	ErrorDescription string `json:"error_description"`
	Wallet           string `json:"wallet"`
}

type TaskRequest struct {
	APIKey string `json:"api_key"`
	TaskId string `json:"taskId"`
}

type TaskResult[TS TaskSolution] struct {
	ErrorId  int    `json:"errorId"`
	Solution TS     `json:"solution"`
	Status   string `json:"status"`
}

type TaskResultRaw struct {
	ErrorId  int             `json:"errorId"`
	Solution json.RawMessage `json:"solution"`
	Status   string          `json:"status"`
}

// https://apidocs.salamoonder.com/tasks/kasada-solver
type KasadaOptions struct {
	Pjs    string `json:"pjs"`
	CdOnly bool   `json:"cdOnly"`
}

// https://apidocs.salamoonder.com/tasks/kasada-solver
type KasadaSolution struct {
	UserAgent string `json:"user-agent"`
	XIsHuman  string `json:"x-is-human"`
	XKpsdkCd  string `json:"x-kpsdk-cd"`
	XKpsdkCr  string `json:"x-kpsdk-cr"`
	XKpsdkCt  string `json:"x-kpsdk-ct"`
	XKpsdkR   string `json:"x-kpsdk-r"`
	XKpsdkSt  string `json:"x-kpsdk-st"`
}

// https://apidocs.salamoonder.com/tasks/incapsula/reese84
type Reese84Options struct {
	Website      string `json:"website"`
	SubmitPayload bool  `json:"submit_payload"`
}

// https://apidocs.salamoonder.com/tasks/incapsula/reese84
type Reese84SubmitPayloadSolution struct {
	Token      string `json:"token"`
	RenewInSec int    `json:"renewInSec"`
	UserAgent  string `json:"user-agent"`
}

// https://apidocs.salamoonder.com/tasks/incapsula/reese84
type Reese84Solution struct {
	Payload        string `json:"payload"`
	UserAgent      string `json:"user-agent"`
	AcceptLanguage string `json:"accept-language"`
}

// https://apidocs.salamoonder.com/tasks/incapsula/utmvc
type UutmvcOptions struct {
	Website string `json:"website"`
}

// https://apidocs.salamoonder.com/tasks/incapsula/utmvc
type UutmvcSolution struct {
	UserAgent string `json:"user-agent"`
	Utmvc     string `json:"utmvc"`
}

// https://apidocs.salamoonder.com/tasks/twitch/scraper
type TwitchScraperOptions struct{}

// https://apidocs.salamoonder.com/tasks/twitch/scraper
type TwitchScraperSolution struct {
	Biography      string `json:"biography"`
	ProfilePicture string `json:"profile_picture"`
	Username       string `json:"username"`
}

// https://apidocs.salamoonder.com/tasks/twitch/integrity
//
// Deprecated: The old "Local Integrity" task has been removed as it’s no longer needed. 
// New structure for solutions is TwitchIntegrityOptions.
//
// https://t.me/salamoonder_telegram/1317
type TwitchPublicIntegrityOptions struct {
	Proxy       string `json:"proxy"`
	AccessToken string `json:"access_token"`
	DeviceID    string `json:"deviceId"`
	ClientID    string `json:"clientId"`
}

// https://apidocs.salamoonder.com/tasks/twitch/integrity
//
// Deprecated: The old "Local Integrity" task has been removed as it’s no longer needed. 
// New structure for solutions is TwitchIntegritySolution.
//
// https://t.me/salamoonder_telegram/1317
type TwitchPublicIntegritySolution struct {
	DeviceID       string `json:"device_id"`
	Proxy          string `json:"proxy"`
	IntegrityToken string `json:"integrity_token"`
	UserAgent      string `json:"user-agent"`
	ClientID       string `json:"client-id"`
}

// https://apidocs.salamoonder.com/tasks/twitch/local-integrity
//
// Deprecated: The old "Local Integrity" task has been removed as it’s no longer needed. 
// Twitch used to require a local integrity token to create an account, 
// but now you can generate accounts using just use KasadaOptions.
//
// https://t.me/salamoonder_telegram/1317
type TwitchLocalIntegrityOptions struct {
	Proxy    string `json:"proxy"`
	DeviceID string `json:"deviceId"`
	ClientID string `json:"clientId"`
}

// https://apidocs.salamoonder.com/tasks/twitch/local-integrity
//
// Deprecated: The old "Local Integrity" task has been removed as it’s no longer needed. 
// Twitch used to require a local integrity token to create an account, 
// but now you can generate accounts using just use KasadaSolution.
//
// https://t.me/salamoonder_telegram/1317
type TwitchLocalIntegritySolution struct {
	DeviceID       string `json:"device_id"`
	IntegrityToken string `json:"integrity_token"`
	Proxy          string `json:"proxy"`
	UserAgent      string `json:"user-agent"`
	ClientID       string `json:"client-id"`
}

// https://apidocs.salamoonder.com/tasks/twitch/integrity
type TwitchIntegrityOptions struct {
	AccessToken string `json:"access_token"`
	DeviceID    string `json:"deviceId"`
	ClientID    string `json:"clientId"`
}

// https://apidocs.salamoonder.com/tasks/twitch/integrity
type TwitchIntegritySolution struct {
	DeviceID       string `json:"device_id"`
	IntegrityToken string `json:"integrity_token"`
	UserAgent      string `json:"user-agent"`
	ClientID       string `json:"client-id"`
}