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
	if !onix.ServiceExists(args.SecondPos) {
		console.Fail(fmt.Sprintf("Service \"%s\" does not exist. It's your word aginst onix.",
			args.SecondPos))
	}

	url := bakeUrl(args.SecondPos, args.Version)

	resp, err := http.Post(url, "application/json", nil)
	if err != nil {
		console.BigFail("Failed to call ditto to bake service")
	}
	defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        console.Fail(fmt.Sprintf("Failed to read ditto response body, that's bad :-(", resp.StatusCode))
    }

	if resp.StatusCode == 200 {
		fmt.Println("Sucessfully baked application:", args.SecondPos,
			"with version:", args.Version)
		fmt.Println(string(body))
	} else {
		fmt.Println("Non 200 response from onix: ", resp.StatusCode)
        console.Fail(fmt.Sprintf("Response body was: %s", body))
	}
}
