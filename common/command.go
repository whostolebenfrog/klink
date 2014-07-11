package common

import (
	"os"
	"regexp"
	"runtime"
)

type Command struct {
	Action      string
	Application string
	Debug       bool
	Description string
	Email       string
	Environment string
	Format      string
	FourthPos   string
	Message     string
	Name        string
	Owner       string
	SecondPos   string
	Silent      bool
	Status      string
	ThirdPos    string
	Type        string
	Value       string
	Version     string
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
