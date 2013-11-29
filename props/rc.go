package props

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	common "nokia.com/klink/common"
	console "nokia.com/klink/console"
    "reflect"
	"os"
)

type RCProps struct {
	Username     string `json:"username"`
    LastUpdated  int32  `json:"lastUpdated"`
	DoctorHasRun string `json:"doctorHasRun"`
}

// Returns the current username
func GetUsername() string {
	EnsureRCFile()
	return GetRCProperties().Username
}

// Creates a klinkrc file and prompts the user for a username
func createRCFile() {
    rcPath := RCFilePath()
	fmt.Println(fmt.Sprintf("\nNo home file found at: %s Creating one for you.\n", rcPath))
	console.Green()
	fmt.Println("Please enter your brislabs username:\n")
	console.Reset()

	var username string
	fmt.Scan(&username)

	rcProps := RCProps{}
    rcProps.Username = username

    writeRCProperties(rcProps)

	fmt.Println(fmt.Sprintf("\nThanks %s, I've created a home file for you.", username))
}

// Ensures that the user has a rc file. If it doesn't exist create it
// after prompting for a username
func EnsureRCFile() {
	if Exists(RCFilePath()) {
		return
	} else {
		createRCFile()
	}
}

// Returns the path of the klink rc file
func RCFilePath() string {
	return common.UserHomeDir() + "/.klinkrc"
}

// Updates the users rc file properties
func UpdateRCProperties(name string, value string) {
    rcProps := GetRCProperties()
    reflect.ValueOf(&rcProps).Elem().FieldByName(name).SetString(value)
    fmt.Println(rcProps)
    writeRCProperties(rcProps)
}

// Returns the users RC properties
func GetRCProperties() RCProps {
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

func writeRCProperties(rcProps RCProps) {
	rcBytes, err := json.Marshal(rcProps)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(RCFilePath(), rcBytes, 0755)
	if err != nil {
		panic(err)
	}
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

func GetLastUpdated() int32 {
    return GetRCProperties().LastUpdated
}

// Write the last time we checked for an update
func SetLastUpdated(t int32) {
    props := GetRCProperties()
    props.LastUpdated = t
    writeRCProperties(props)
}
