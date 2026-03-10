package salamoonder

import (
	"encoding/json"
	"reflect"
)

func init() {
	registerAllowedTypes(
		KasadaStandardOptions{},
		KasadaPayloadOptions{},
		AkamaiWebOptions{},
		AkamaiSBSDOptions{},
		Reese84Options{},
		UutmvcOptions{},
		DataDomeInterstitialOptions{},
		DataDomeSliderOptions{},
		TwitchScraperOptions{},
		TwitchIntegrityOptions{},
	)
}

func registerAllowedTypes(values ...any) {
	for _, v := range values {
		allowedTaskTypes = append(allowedTaskTypes, reflect.TypeOf(v))
	}
}

type (
	TaskOptions interface {
		KasadaStandardOptions | KasadaPayloadOptions | AkamaiWebOptions | AkamaiSBSDOptions | Reese84Options | UutmvcOptions | DataDomeInterstitialOptions | DataDomeSliderOptions | TwitchScraperOptions | TwitchIntegrityOptions
	}

	TaskSolution interface {
		KasadaStandardSolution | KasadaPayloadSolution | AkamaiWebSolution | AkamaiSBSDSolution | Reese84Solution | Reese84SubmitPayloadSolution | UutmvcSolution | DataDomeInterstitialSolution | DataDomeSliderSolution | TwitchScraperSolution | TwitchIntegritySolution
	}

	CreateTaskRequest struct {
		ApiKey string      `json:"api_key"`
		Task   interface{} `json:"task"`
	}

	CreateTaskResult struct {
		ErrorCode        int    `json:"error_code"`
		ErrorDescription string `json:"error_description"`
		TaskId           string `json:"taskId"`
	}

	// https://apidocs.salamoonder.com/endpoint/getBalance
	CreateTaskBalanceResult struct {
		ErrorCode        int    `json:"error_code"`
		ErrorDescription string `json:"error_description"`
		Wallet           string `json:"wallet"`
	}

	TaskRequest struct {
		APIKey string `json:"api_key"`
		TaskId string `json:"taskId"`
	}

	TaskResult[TS TaskSolution] struct {
		ErrorId  int    `json:"errorId"`
		Solution TS     `json:"solution"`
		Status   string `json:"status"`
	}

	TaskResultRaw struct {
		ErrorId  int             `json:"errorId"`
		Solution json.RawMessage `json:"solution"`
		Status   string          `json:"status"`
	}

	// https://apidocs.salamoonder.com/tasks/kasada/standard
	KasadaStandardOptions struct {
		Pjs    string `json:"pjs"`
		CdOnly bool   `json:"cdOnly"`
	}

	// https://apidocs.salamoonder.com/tasks/kasada/standard
	KasadaStandardSolution struct {
		UserAgent string `json:"user-agent"`
		XIsHuman  string `json:"x-is-human"`
		XKpsdkCd  string `json:"x-kpsdk-cd"`
		XKpsdkCr  string `json:"x-kpsdk-cr"`
		XKpsdkCt  string `json:"x-kpsdk-ct"`
		XKpsdkR   string `json:"x-kpsdk-r"`
		XKpsdkSt  string `json:"x-kpsdk-st"`
	}

	// https://apidocs.salamoonder.com/tasks/kasada/payload
	KasadaPayloadOptions struct {
		URL           string `json:"url"`
		ScriptURL     string `json:"script_url"`
		ScriptContent string `json:"script_content"`
	}

	// https://apidocs.salamoonder.com/tasks/kasada/payload
	KasadaPayloadSolution struct {
		Headers struct {
			XKpsdkCt string `json:"x-kpsdk-ct"`
			XKpsdkDt string `json:"x-kpsdk-dt"`
			XKpsdkIm string `json:"x-kpsdk-im"`
			XKpsdkV  string `json:"x-kpsdk-v"`
		} `json:"headers"`
		Payload   string `json:"payload"`
		UserAgent string `json:"user-agent"`
	}

	// https://apidocs.salamoonder.com/tasks/akamai/web
	AkamaiWebOptions struct {
		Type      string `json:"type"`
		URL       string `json:"url"`
		Abck      string `json:"abck"`
		Bmsz      string `json:"bmsz"`
		Script    string `json:"script"`
		SensorUrl string `json:"sensor_url"`
		Count     int64  `json:"count"`
		Data      string `json:"data"`
		UserAgent string `json:"user_agent"`
	}

	// https://apidocs.salamoonder.com/tasks/akamai/web
	AkamaiWebSolution struct {
		Payload   map[string]any `json:"payload"`
		Data      map[string]any `json:"data"`
		UserAgent string         `json:"user-agent"`
	}

	// https://apidocs.salamoonder.com/tasks/akamai/sbsd
	AkamaiSBSDOptions struct {
		URL       string `json:"url"`
		Cookie    string `json:"cookie"`
		SbsdURL   string `json:"sbsd_url"`
		Script    string `json:"script"`
		UserAgent string `json:"user_agent"`
	}

	// https://apidocs.salamoonder.com/tasks/akamai/sbsd
	AkamaiSBSDSolution struct {
		Payload   string `json:"payload"`
		UserAgent string `json:"user-agent"`
	}

	// https://apidocs.salamoonder.com/api-documentation/tasks/incapsula/reese84
	Reese84Options struct {
		Website       string `json:"website"`
		SubmitPayload bool   `json:"submit_payload"`
	}

	// https://apidocs.salamoonder.com/tasks/incapsula/reese84
	Reese84SubmitPayloadSolution struct {
		Token      string `json:"token"`
		RenewInSec int    `json:"renewInSec"`
		UserAgent  string `json:"user-agent"`
	}

	// https://apidocs.salamoonder.com/tasks/incapsula/reese84
	Reese84Solution struct {
		Payload        string `json:"payload"`
		UserAgent      string `json:"user-agent"`
		AcceptLanguage string `json:"accept-language"`
	}

	// https://apidocs.salamoonder.com/tasks/incapsula/utmvc
	UutmvcOptions struct {
		Website string `json:"website"`
	}

	// https://apidocs.salamoonder.com/tasks/incapsula/utmvc
	UutmvcSolution struct {
		UserAgent string `json:"user-agent"`
		Utmvc     string `json:"utmvc"`
	}

	// https://apidocs.salamoonder.com/tasks/twitch/scraper
	TwitchScraperOptions struct{}

	// https://apidocs.salamoonder.com/tasks/twitch/scraper
	TwitchScraperSolution struct {
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
	TwitchPublicIntegrityOptions struct {
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
	TwitchPublicIntegritySolution struct {
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
	TwitchLocalIntegrityOptions struct {
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
	TwitchLocalIntegritySolution struct {
		DeviceID       string `json:"device_id"`
		IntegrityToken string `json:"integrity_token"`
		Proxy          string `json:"proxy"`
		UserAgent      string `json:"user-agent"`
		ClientID       string `json:"client-id"`
	}

	// https://apidocs.salamoonder.com/tasks/twitch/integrity
	TwitchIntegrityOptions struct {
		AccessToken string `json:"access_token"`
		DeviceID    string `json:"deviceId"`
		ClientID    string `json:"clientId"`
	}

	// https://apidocs.salamoonder.com/tasks/twitch/integrity
	TwitchIntegritySolution struct {
		DeviceID       string `json:"device_id"`
		IntegrityToken string `json:"integrity_token"`
		UserAgent      string `json:"user-agent"`
		ClientID       string `json:"client-id"`
	}

	// https://apidocs.salamoonder.com/tasks/datadome/interstitial
	DataDomeInterstitialOptions struct {
		CaptchaURL  string `json:"captcha_url"`
		UserAgent   string `json:"user_agent"`
		CountryCode string `json:"country_code"`
	}

	// https://apidocs.salamoonder.com/tasks/datadome/interstitial
	DataDomeInterstitialSolution struct {
		Cookie    string `json:"cookie"`
		UserAgent string `json:"user-agent"`
	}

	// https://apidocs.salamoonder.com/tasks/datadome/slider
	DataDomeSliderOptions struct {
		CaptchaURL  string `json:"captcha_url"`
		UserAgent   string `json:"user_agent"`
		CountryCode string `json:"country_code"`
	}

	// https://apidocs.salamoonder.com/tasks/datadome/slider
	DataDomeSliderSolution struct {
		Cookie    string `json:"cookie"`
		UserAgent string `json:"user-agent"`
	}
)
