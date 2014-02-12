package ditto

import (
	"fmt"
	"io"
	"net/http"
	common "nokia.com/klink/common"
	console "nokia.com/klink/console"
	onix "nokia.com/klink/onix"
	"os"
)

func dittoUrl(end string) string {
	return "http://ditto.brislabs.com:8080/1.x" + end
}

func bakeUrl(app string, version string) string {
	return fmt.Sprintf(dittoUrl("/bake/%s/%s"), app, version)
}

func AllowProd(args common.Command) {
	if args.SecondPos == "" {
		console.Fail("Application must be provided as the second positional argument")
	}

	url := dittoUrl("/make-public/" + args.SecondPos)

	resp, err := http.Post(url, "application/json", nil)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode == 200 {
		fmt.Println("Success")
	} else {
		panic(fmt.Sprintf("%d response calling URL: ", resp.StatusCode))
	}
}

func DoBake(url string) {
	httpClient := common.NewTimeoutClient()
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Println(err)
		console.Fail(fmt.Sprintf("Failed to make a request to: %s", url))
	}
	defer resp.Body.Close()

	if resp.StatusCode == 503 {
		console.Fail("Ditto is currently not available. This is most likely due it being redeployed. If it's not back in 10 minutes ask in Campfire or speak to Ben Griffiths.")
	} else if resp.StatusCode == 404 {
		console.Fail("Sorry, the RPM for this application is not yet available. Wait a few minutes and then try again.")
	} else if resp.StatusCode != 200 {
		fmt.Println(fmt.Sprintf("Got %d response calling ditto to bake ami.", resp.StatusCode))
		io.Copy(os.Stdout, resp.Body)
		panic("\nFailed to bake ami.")
	}

	io.Copy(os.Stdout, resp.Body)
}

// Bake the ami
func Bake(args common.Command) {
	if args.SecondPos == "" {
		console.Fail("Application must be supplied as second positional argument")
	}
	if args.Version == "" {
		console.Fail("Version must be supplied using --version")
	}
	if !onix.AppExists(args.SecondPos) {
		console.Fail(fmt.Sprintf("Application '%s' does not exist. It's your word against onix!",
			args.SecondPos))
	}

	url := bakeUrl(args.SecondPos, args.Version)
    DoBake(url)
}

type Ami struct {
	Name    string
	ImageId string
}

// FindAmis using the service name for the second positional command line arg
// Prints out a list of the most recent ami names and image ids
func FindAmis(args common.Command) {
	if args.SecondPos == "" {
		console.Fail("Application must be supplied as second positional argument")
	}

	amis := make([]Ami, 10)
	common.GetJson(dittoUrl(fmt.Sprintf("/amis/%s", args.SecondPos)), &amis)

	for key := range amis {
		fmt.Print(amis[key].Name, " : ")
		console.Brown()
		fmt.Print(amis[key].ImageId)
		console.Grey()
		fmt.Println()
	}
	console.Reset()
}

// ditto helps to lock, unlock and clean amis
// not intended to be a part of the public klink functionality
func Helpers(args common.Command) {
    command := args.SecondPos
	if command == "lock" {
		lockUrl := dittoUrl("/lock")
		fmt.Println(common.PostJson(lockUrl, nil))
	} else if command == "unlock" {
		unlockUrl := dittoUrl("/unlock")
		fmt.Println(common.PostJson(unlockUrl, nil))
	} else if command == "clean" {
		cleanUrl := dittoUrl("/clean")
		if args.ThirdPos == "" {
			cleanUrl += "all"
		} else {
			cleanUrl += args.ThirdPos
		}
		fmt.Println(common.PostJson(cleanUrl, nil))
	} else if command == "entertainment" {
        bakeUrl := dittoUrl("/bake/entertainment-ami")
        DoBake(bakeUrl)
    } else if command == "public" {
        bakeUrl := dittoUrl("/bake/public-ami")
        DoBake(bakeUrl)
    } else if command == "inprogress" {
        progUrl := dittoUrl("/inprogress")
        fmt.Println(common.GetString(progUrl))
    } else {
        console.Fail("Requires a second arg: lock, unlock, clean, entertainment, public or inprogress")
    }
}
