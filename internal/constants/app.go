package constants

import "time"

var (
	AppName        = "app"
	Version        = "0.1.0"
	BuildTimestamp = time.Now().UTC().Format("2006-01-02T15:04:05")
)
