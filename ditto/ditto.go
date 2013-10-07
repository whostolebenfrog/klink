package ditto

import (
	"fmt"
	"io/ioutil"
	"net/http"
	common "nokia.com/klink/common"
	console "nokia.com/klink/console"
    onix "nokia.com/klink/onix"
)

func dittoUrl(end string) string {
	return "http://ditto.brislabs.com:8080/1.x" + end
}

func bakeUrl(app string, version string) string {
	return fmt.Sprintf(dittoUrl("/bake/%s/%s"), app, version)
}

func Bake(args common.Command) {
	if args.SecondPos == "" {
		console.Fail("Application must be supplied as second positional argument")
	}
	if args.Version == "" {
		console.Fail("Version must be supplied using --version")
	}
	if !onix.AppExists(args.SecondPos) {
		console.Fail(fmt.Sprintf("Application '%s' does not exist. It's your word aginst onix.",
			args.SecondPos))
	}

	url := bakeUrl(args.SecondPos, args.Version)

	resp, err := http.Post(url, "application/json", nil)
	if err != nil {
		console.BigFail("Failed to call ditto to bake application")
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			console.BigFail("Failed to read ditto response body, that's bad :-(")
		}
		fmt.Println("Sucessfully baked application:", args.SecondPos,
			"with version:", args.Version)
		fmt.Println(string(body))
	} else if resp.StatusCode == 404 {
		fmt.Println("Sorry, the RPM for this application is not yet available. Wait a few minutes and then try again.")
	}
}
