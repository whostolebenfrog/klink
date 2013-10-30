package props

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	console "nokia.com/klink/console"
	"os"
	"runtime"
)

type RCProps struct {
	Username string `json:"username"`
}

// Returns the current username
func GetUsername() string {
	EnsureRCFile()
	return getRCProperties().Username
}

// Creates a klinkrc file and prompts the user for a username
func createRCFile(rcPath string) {
	fmt.Println(fmt.Sprintf("\nNo home file found at: %s Creating one for you.\n", rcPath))
	console.Green()
	fmt.Println("Please enter your brislabs username:\n")
	console.Reset()

	var username string
	fmt.Scan(&username)

	rcProps := RCProps{username}

	rcBytes, err := json.Marshal(rcProps)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(rcPath, rcBytes, 0755)
	if err != nil {
		panic(err)
	}

	fmt.Println(fmt.Sprintf("\nThanks %s, I've created a home file for you.", username))
}

// Ensures that the user has a rc file. If it doesn't exist create it
// after prompting for a username
func EnsureRCFile() {
	if Exists(RCFilePath()) {
		return
	} else {
		createRCFile(RCFilePath())
	}
}

// Returns the path of the klink rc file
func RCFilePath() string {
	return userHomeDir() + "/.klinkrc"
}

// Updates the users rc file properties
func updateRCProperties() {
}

// Returns the users RC properties
func getRCProperties() RCProps {
	rcBytes, err := ioutil.ReadFile(RCFilePath())

	if err != nil {
		panic(err)
	}

	rcProps := RCProps{}
	err = json.Unmarshal(rcBytes, &rcProps)

	if err != nil {
		panic(err)
	}
	return rcProps
}

// Exists reports whether the named file or directory exists.
func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// This is required as user.Current fails on darwin when cross compiled from linux.
// If anyone reading this understands enough about builders to fix it - this seems
// to be the same issue:
// https://groups.google.com/forum/#!topic/golang-dev/zzBrnKMYctQ
func userHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}
