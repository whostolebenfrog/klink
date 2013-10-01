package exploud

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	common "nokia.com/klink/common"
	console "nokia.com/klink/console"
)

type DeployRequest struct {
	Ami         string `json:"ami"`
	Environment string `json:"environment"`
}

// TODO: handle message when exploud supports it
func Exploud(args common.Command) {
    if args.SecondPos == "" {
        console.Fail("Must supply an application name as the second positional argument")
    }
    if args.Ami == "" {
        console.Fail("Must supply an ami to deploy using --ami")
    }

	deployUrl := fmt.Sprintf("http://exploud.brislabs.com:8080/1.x/applications/%s/deploy", args.SecondPos)

    // TODO: use http lib and add polling here!
	deployRequest := DeployRequest{args.Ami, "dev"}
	b, err := json.Marshal(deployRequest)
	if err != nil {
		console.BigFail("Unable to create exploud deploy requset body")
	}
	fmt.Println("Calling exploud:", deployUrl, string(b))

	resp, err := http.Post(deployUrl, "application/json", bytes.NewReader(b))
	if err != nil {
		console.BigFail("Error response from exploud")
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			console.BigFail("Failed to read exploud response body, that's bad :-(")
		}
		fmt.Println("200 response from exploud - that's good!")
		fmt.Println(string(body))
	}
}
