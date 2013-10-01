package ditto

import (
	"fmt"
	"io/ioutil"
	"net/http"
	common "nokia.com/klink/common"
	console "nokia.com/klink/console"
)

func dittoUrl(end string) string {
    return "http://localhost:8080/1.x" + end
}

func bakeUrl(app string, version string) {
    return fmt.Sprintf(dittoUrl("/bake/%s/%s"), app, version)
}

func Bake(command common.Command) {
	if command.Version == "" {
		console.Fail("Args version must be supplied")
	}
	if command.Application == "" {
		console.Fail("Application must be supplied")
	}

    url := bakeUrl(command.Application, command.Version)

	resp, err := http.Post(url, "application/json", nil)
	if err != nil {
		console.BigFail("Failed to call ditto to bake service")
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			console.BigFail("Failed to read ditto response body, that's bad :-(")
		}
		fmt.Println("Sucessfully baked application:", command.Application, "with version:", command.Version)
		fmt.Println(string(body))
	}
}
