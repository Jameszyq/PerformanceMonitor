package model

var Config ConfigInfo

type ConfigInfo struct {
	WebHookAddress         string
	CollectionIntervalCorn string
	AlarmCount             int
	ServerAlias            string
	ServerIp               string
	CollectionIntervalUnit string
	WarnIndex              float64
}
