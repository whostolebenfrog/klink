package ditto

import (
	"fmt"
	"io"
	"net/http"
	common "nokia.com/klink/common"
	console "nokia.com/klink/console"
	onix "nokia.com/klink/onix"
	"os"
    "time"
	"regexp"
)

func Init() {
	common.Register(
		common.Component{"ditto", Helpers,
			"Various helpers; lock, unlock, clean, build public and ent base amis"},
		common.Component{"bake", Bake,
			"{app} -v {version} Bakes an AMI for {app} with version {version}"},
		common.Component{"allow-prod", AllowProd,
			"{app} Allows the prod aws account access to the supplied application"},
		common.Component{"list-amis", FindAmis,
			"{app} Lists the latest amis for the supplied application name"},
		common.Component{"latest-bake", LatestBake,
			"{app} Outputs the latest baked version of the specified application"})
}

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
	httpClient := common.NewTimeoutClient(5*time.Second, 1200*time.Second)
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
		fmt.Printf("Got %d response calling ditto to bake ami.\n", resp.StatusCode)
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
	application := args.SecondPos

	if application == "" {
		console.Fail("Application must be supplied as second positional argument")
	}

	amis := make([]Ami, 10)
	common.GetJson(dittoUrl(fmt.Sprintf("/amis/%s", application)), &amis)

	for key := range amis {
		fmt.Print(amis[key].Name, " : ")
		console.Brown()
		fmt.Print(amis[key].ImageId)
		console.Grey()
		fmt.Println()
	}
	console.Reset()
}

func LatestBake(args common.Command) {
	application := args.SecondPos

	if application == "" {
		console.Fail("Application must be supplied as second positional argument")
	}

	amis := make([]Ami, 10)
	common.GetJson(dittoUrl(fmt.Sprintf("/amis/%s", application)), &amis)

	version := parseVersionFrom(amis[0])

	fmt.Print(version)
	fmt.Println()

	console.Reset()
}

func parseVersionFrom(ami Ami) string {
	versionRegexp := regexp.MustCompile("[0-9.]+")

	return versionRegexp.FindString(ami.Name)
}

// ditto helps to lock, unlock and clean amis
// not intended to be a part of the public klink functionality
func Helpers(args common.Command) {
	switch args.SecondPos {
	case "lock":
		lockUrl := dittoUrl("/lock")
		fmt.Println(common.PostJson(lockUrl, nil))
	case "unlock":
		unlockUrl := dittoUrl("/unlock")
		fmt.Println(common.PostJson(unlockUrl, nil))
	case "clean":
		cleanUrl := dittoUrl("/clean/")
		if args.ThirdPos == "" {
			cleanUrl += "all"
		} else {
			cleanUrl += args.ThirdPos
		}
		fmt.Println(common.PostJson(cleanUrl, nil))
	case "entertainment":
		bakeUrl := dittoUrl("/bake/entertainment-ami")
		DoBake(bakeUrl)
	case "public":
		bakeUrl := dittoUrl("/bake/public-ami")
		DoBake(bakeUrl)
	case "inprogress":
		progUrl := dittoUrl("/inprogress")
		fmt.Println(common.GetString(progUrl))
	default:
		console.Fail("Requires a second arg: lock, unlock, clean, entertainment, public or inprogress")
	}
}
