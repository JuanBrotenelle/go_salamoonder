package salamoonder

import "encoding/json"

type TaskOptions interface {
	KasadaOptions | Reese84Options | UutmvcOptions | TwitchScraperOptions | TwitchPublicIntegrityOptions | TwitchLocalIntegrityOptions
}

type TaskSolution interface {
	KasadaSolution | Reese84Solution | Reese84SubmitPayloadSolution | UutmvcSolution | TwitchScraperSolution | TwitchPublicIntegritySolution | TwitchLocalIntegritySolution
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

type KasadaOptions struct {
	Pjs    string `json:"pjs"`
	CdOnly bool   `json:"cdOnly"`
}

type KasadaSolution struct {
	UserAgent string `json:"user-agent"`
	XIsHuman  string `json:"x-is-human"`
	XKpsdkCd  string `json:"x-kpsdk-cd"`
	XKpsdkCr  string `json:"x-kpsdk-cr"`
	XKpsdkCt  string `json:"x-kpsdk-ct"`
	XKpsdkR   string `json:"x-kpsdk-r"`
	XKpsdkSt  string `json:"x-kpsdk-st"`
}

type Reese84Options struct {
	Website      string `json:"website"`
	SubmitPayload bool  `json:"submit_payload"`
}

type Reese84SubmitPayloadSolution struct {
	Token      string `json:"token"`
	RenewInSec int    `json:"renewInSec"`
	UserAgent  string `json:"user-agent"`
}

type Reese84Solution struct {
	Payload        string `json:"payload"`
	UserAgent      string `json:"user-agent"`
	AcceptLanguage string `json:"accept-language"`
}

type UutmvcOptions struct {
	Website string `json:"website"`
}

type UutmvcSolution struct {
	UserAgent string `json:"user-agent"`
	Utmvc     string `json:"utmvc"`
}

type TwitchScraperOptions struct{}

type TwitchScraperSolution struct {
	Biography      string `json:"biography"`
	ProfilePicture string `json:"profile_picture"`
	Username       string `json:"username"`
}

type TwitchPublicIntegrityOptions struct {
	Proxy       string `json:"proxy"`
	AccessToken string `json:"access_token"`
	DeviceID    string `json:"deviceId"`
	ClientID    string `json:"clientId"`
}

type TwitchPublicIntegritySolution struct {
	DeviceID       string `json:"device_id"`
	IntegrityToken string `json:"integrity_token"`
	Proxy          string `json:"proxy"`
	UserAgent      string `json:"user-agent"`
	ClientID       string `json:"client-id"`
}

type TwitchLocalIntegrityOptions struct {
	Proxy    string `json:"proxy"`
	DeviceID string `json:"deviceId"`
	ClientID string `json:"clientId"`
}

type TwitchLocalIntegritySolution struct {
	DeviceID       string `json:"device_id"`
	IntegrityToken string `json:"integrity_token"`
	Proxy          string `json:"proxy"`
	UserAgent      string `json:"user-agent"`
	ClientID       string `json:"client-id"`
}
