package ditto

import (
	"fmt"
	"io"
	"net/http"
	common "nokia.com/klink/common"
	console "nokia.com/klink/console"
	exploud "nokia.com/klink/exploud"
	"os"
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
	if !exploud.AppExists(args.SecondPos) {
		console.Fail(fmt.Sprintf("Application '%s' does not exist. It's your word aginst exploud!",
			args.SecondPos))
	}

	url := bakeUrl(args.SecondPos, args.Version)

    httpClient := common.NewTimeoutClient()
    req, _ := http.NewRequest("POST", url, nil)
    req.Header.Add("Content-Type", "application/json")

    resp, err := httpClient.Do(req)
    if err != nil {
        fmt.Println(err)
        console.Fail(fmt.Sprintf("Failed to make a request to: %s", url))
    }
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		console.Fail("Sorry, the RPM for this application is not yet available. Wait a few minutes and then try again.")
	} else if resp.StatusCode != 200 {
		console.Fail(fmt.Sprintf("Non 200 response from ditto: ", resp.StatusCode))
	}
	defer resp.Body.Close()

	io.Copy(os.Stdout, resp.Body)
}

type Ami struct {
    Name string
    ImageId string
}

func FindAmis(args common.Command) {
    if args.SecondPos == "" {
        console.Fail("Application must be supplied as second positional argument")
    }

    amis := make([]Ami, 10)
    err := common.GetJson(dittoUrl(fmt.Sprintf("/amis/%s", args.SecondPos)), &amis)

    if err != nil {
        fmt.Println(err)
        console.Fail("Could not list amis from ditto")
    }

    for key := range amis {
        fmt.Println(fmt.Sprintf("%s : \033[32m%s\033[37m", amis[key].Name, amis[key].ImageId))
    }
    fmt.Println("")
}
