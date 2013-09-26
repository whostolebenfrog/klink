package ditto

import (
	"fmt"
	"io/ioutil"
	"net/http"
	args "nokia.com/klink/args"
	console "nokia.com/klink/console"
)

func Bake(command args.Command) {
	if command.Version == "" {
		console.Fail("Args version must be supplied")
	}
	if command.Application == "" {
		console.Fail("Application must be supplied")
	}
	bakeUrl := fmt.Sprintf("http://localhost:8080/1.x/bake/%s/%s", command.Application, command.Version)
	resp, err := http.Post(bakeUrl, "application/json", nil)
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
