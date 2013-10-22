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
	ThirdPos    string
	FourthPos   string
	Version     string
	Description string
	Email       string
	Owner       string
	Silent      bool
	Debug       bool
	Name        string
	Value       string
}

// Do we have to do stupid shit to get around windows being a moron?
func IsWindows() bool {
	matched, _ := regexp.MatchString(".*windows.*", runtime.GOOS)
	return matched
}
