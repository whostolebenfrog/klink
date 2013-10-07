package exploud

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	common "nokia.com/klink/common"
	console "nokia.com/klink/console"
	onix "nokia.com/klink/onix"
)

type DeployRequest struct {
	Ami         string `json:"ami"`
	Environment string `json:"environment"`
}

type CreateAppRequest struct {
        Description string `json:"description"`
	Email string `json:"email"`
        Owner string `json:"owner"`
}

func exploudUrl(end string) string {
    return "http://exploud.brislabs.com:8080/1.x" + end
}

// TODO: use http lib and add polling here!
// TODO: handle message when exploud supports it
func Exploud(args common.Command) {
	if args.SecondPos == "" {
		console.Fail("Must supply an application name as the second positional argument")
	}
	if args.Ami == "" {
		console.Fail("Must supply an ami to deploy using --ami")
	}
	if !onix.AppExists(args.SecondPos) {
		console.Fail(fmt.Sprintf("Application \"%s\" does not exist. It's your word aginst onix.",
			args.SecondPos))
	}

	deployUrl := fmt.Sprintf(exploudUrl("/applications/%s/deploy"), args.SecondPos)

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

func CreateApp(args common.Command) {
	if args.SecondPos == "" || strings.Index(args.SecondPos, "-") == 0 {
		console.Fail("Must supply an application name as second positional argument")
	}

        if args.Description == "" || args.Email == "" || args.Owner == "" {
		console.Fail("Don't be lazy! You must supply owner, email and description values")
	}

	fmt.Println(fmt.Sprintf("Calling exploud to create application %s with description: %s, email: %s, owner: %s",
		args.SecondPos, args.Description, args.Email, args.Owner))

	createBody := CreateAppRequest{args.Description, args.Email, args.Owner}

	response, err := common.PutJson(exploudUrl("/applications/" + args.SecondPos), createBody)

	if err != nil {
		fmt.Println(err)
		console.BigFail("Unable to register new application with exploud")
	}

	fmt.Println("Exploud has created our application for us!")
	fmt.Println(response)
}
