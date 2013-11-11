package common

import (
	"os"
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

// This is required as user.Current fails on darwin when cross compiled from linux.
// If anyone reading this understands enough about builders to fix it - this seems
// to be the same issue:
// https://groups.google.com/forum/#!topic/golang-dev/zzBrnKMYctQ
func UserHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}
