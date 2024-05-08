package models

import (
	"time"
)

type LogModel struct {
	LogId string `json:"log_id"`
	Timestamp time.Time `json:"timestamp"`
	UserId int `json:"user_id"`
	Actions string `json:"action"`
	Resource string `json:"resource"`
	Status string `json:"status"`
	ClientIp string `json:"client_ip"`
	ClientDevice string `json:"client_device"`
	ClientOs string `json:"client_os"`
	ClientOsVersion string `json:"client_os_ver"`
	ClientBrowser string `json:"client_browser"`
	ClientBrowserVersion string `json:"client_browser_ver"`
	Duration time.Duration `json:"duration"`
	Errors  string `json:"errors"`
}