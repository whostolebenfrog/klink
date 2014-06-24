package props

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	common "nokia.com/klink/common"
	console "nokia.com/klink/console"
	"reflect"
)

type RCProps struct {
	Username     string   `json:"username"`
	LastUpdated  int32    `json:"lastUpdated"`
	DoctorHasRun string   `json:"doctorHasRun"`
	SSHUsername  string   `json:"sshUsername"`
	Environments []string `json:"environments"`
}

//////////////////////
// Helper functions //
//////////////////////

// Returns the path of the klink rc file
func rcFilePath() string {
	return common.UserHomeDir() + "/.klinkrc"
}

// write the properties to disk
func writeRCProperties(rcProps RCProps) {
	rcBytes, err := json.Marshal(rcProps)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(rcFilePath(), rcBytes, 0755)
	if err != nil {
		panic(err)
	}
}

// Creates a klinkrc file and prompts the user for a username
func createRCFile() {
	rcPath := rcFilePath()
	fmt.Printf("\nNo home file found at: %s Creating one for you.\n\n", rcPath)
	console.Green()
	fmt.Println("Please enter your brislabs username:\n")
	console.Reset()

	var username string
	fmt.Scan(&username)

	rcProps := RCProps{}
	rcProps.Username = username

	writeRCProperties(rcProps)

	fmt.Printf("\nThanks %s, I've created a home file for you.\n", username)
}

//////////////////////
// public functions //
//////////////////////

// Ensures that the user has a rc file. If it doesn't exist create it
// after prompting for a username
func EnsureRCFile() {
	if common.Exists(rcFilePath()) {
		return
	} else {
		createRCFile()
	}
}

// Returns the users RC properties
func GetRCProperties() RCProps {
	rcBytes, err := ioutil.ReadFile(rcFilePath())

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

// Set the property key / value pair in the users rc file
func Set(name string, value string) {
	rcProps := GetRCProperties()
	reflect.ValueOf(&rcProps).Elem().FieldByName(name).SetString(value)
	writeRCProperties(rcProps)
}

// Get a property by name
func Get(name string) string {
	rcProps := GetRCProperties()
	return reflect.ValueOf(rcProps).FieldByName(name).String()
}

// Return the last time we checked for an update
func GetLastUpdated() int32 {
	return GetRCProperties().LastUpdated
}

// Write the last time we checked for an update
func SetLastUpdated(t int32) {
	props := GetRCProperties()
	props.LastUpdated = t
	writeRCProperties(props)
}

// Returns the list of known environments
func GetEnvironments() []string {
	return GetRCProperties().Environments
}

// Updates the list of known environments
func SetEnvironments(environments []string) {
	props := GetRCProperties()
	props.Environments = environments
	writeRCProperties(props)
}
