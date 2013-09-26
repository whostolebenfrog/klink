package exploud

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	args "nokia.com/klink/args"
	console "nokia.com/klink/console"
)

type DeployRequest struct {
	Ami         string `json:"ami"`
	Environment string `json:"environment"`
}

func Exploud(args args.Command) {
	deployUrl := fmt.Sprintf("http://10.216.138.6:8080/1.x/applications/%s/deploy", args.Application)

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
