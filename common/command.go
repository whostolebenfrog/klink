package common

import (
	"regexp"
	"runtime"
)

type Command struct {
	Action      string
	Ami         string
	Application string
	Environment string
	Message     string
	SecondPos   string
	Version     string
	Description string
	Email       string
	Owner       string
}

// Do we have to do stupid shit to get around windows being a moron?
func IsWindows() bool {
	matched, _ := regexp.MatchString(".*windows.*", runtime.GOOS)
	return matched
}
