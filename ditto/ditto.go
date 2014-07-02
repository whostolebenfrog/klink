package ditto

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"time"

	common "nokia.com/klink/common"
	console "nokia.com/klink/console"
	onix "nokia.com/klink/onix"
)

func Init() {
	common.Register(
		common.Component{"ditto", Helpers,
			"Various helpers; lock, unlock, clean, build public and ent base amis"},
		common.Component{"bake", Bake,
			"{app} [{version}] [-v {version}] Bakes an AMI for {app} with version {version}"},
		common.Component{"allow-prod", AllowProd,
			"{app} Allows the prod aws account access to the supplied application"},
		common.Component{"amis", FindAmis,
			"{app} Lists the latest amis for the supplied application name"},
		common.Component{"images", FindAmis,
			"{app} Lists the latest images for the supplied application name"},
		common.Component{"latest-bake", LatestBake,
			"{app} Outputs the latest baked version of the specified application"},
        common.Component{"delete-ami", DeleteAmi,
            "{service} {ami} Removes the supplied ami, makes it undeployable."})
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

func DoBake(url string, retries int) {
	httpClient := common.NewTimeoutClient(10*time.Second, 2000*time.Second)
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
		if retries > 0 {
			fmt.Println("RPM isn't yet available, retrying in 5 seconds")
			time.Sleep(5 * time.Second)
			DoBake(url, retries-1)
		} else {
			console.Fail("RPM doesn't appear to be available after retrying. Wait a bit longer or check your version is correct / yum repo has it available")
		}
	} else if resp.StatusCode != 200 {
		fmt.Printf("Got %d response calling ditto to bake ami.\n", resp.StatusCode)
		io.Copy(os.Stdout, resp.Body)
		panic("\nFailed to bake ami.")
	} else {
		io.Copy(os.Stdout, resp.Body)
	}
}

// Bake the ami
func Bake(args common.Command) {
	app := args.SecondPos
	if app == "" {
		console.Fail("Application must be supplied as the second argument")
	}

	var version string
	if args.ThirdPos != "" {
		version = args.ThirdPos
	} else {
		version = args.Version
	}
	if version == "" {
		console.Fail("Version must be supplied as the third argument or using --version")
	}

	if !onix.AppExists(app) {
		console.Fail(fmt.Sprintf("Application '%s' does not exist. It's your word against onix!",
			app))
	}

	url := bakeUrl(app, version)
	DoBake(url, 120)
}

type Ami struct {
	Name    string
	ImageId string
	Version string
}

// FindAmis using the application name for the second positional command line arg
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

func LatestAmiFor(application string) Ami {
	amis := make([]Ami, 10)
	common.GetJson(dittoUrl(fmt.Sprintf("/amis/%s", application)), &amis)

	latestAmi := amis[0]
	latestAmi.Version = parseVersionFrom(latestAmi)

	return latestAmi
}

func LatestBake(args common.Command) {
	application := args.SecondPos

	if application == "" {
		console.Fail("Application must be supplied as second positional argument")
	}

	fmt.Print(LatestAmiFor(application).Version)
	fmt.Println()

	console.Reset()
}

func parseVersionFrom(ami Ami) string {
	versionRegexp := regexp.MustCompile("[0-9.]+")

	return versionRegexp.FindString(ami.Name)
}

type Lock struct {
	Message string `json:"message"`
}

// ditto helps to lock, unlock and clean amis
// not intended to be a part of the public klink functionality
func Helpers(args common.Command) {
	switch args.SecondPos {
	case "lock":
		if args.ThirdPos == "" {
			console.Fail("Pass a message you fool.")
		}
		lock := Lock{args.ThirdPos}
		lockUrl := dittoUrl("/lock")
		fmt.Println(common.PostJson(lockUrl, lock))
	case "unlock":
		unlockUrl := dittoUrl("/unlock")
		common.Delete(unlockUrl)
		console.Green()
		fmt.Println("unlock")
		console.Reset()
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
		DoBake(bakeUrl, 120)
	case "public":
		bakeUrl := dittoUrl("/bake/public-ami")
		DoBake(bakeUrl, 120)
	case "inprogress":
		progUrl := dittoUrl("/inprogress")
		fmt.Println(common.GetString(progUrl))
	default:
		console.Fail("Requires a second arg: lock, unlock, clean, entertainment, public or inprogress")
	}
}

// Remove the supplied ami
func DeleteAmi(args common.Command) {
    service := args.SecondPos
    ami := args.ThirdPos

    if service == "" {
        console.Fail("You must supply a service or klink... well you don't want to know.")
    }
    if ami == "" {
        console.Fail("This isn't going to work without an ami now is it?")
    }
    common.FailIfNotAmi(ami)

    common.Delete(dittoUrl(fmt.Sprintf("/%s/amis/%s", service, ami)))
    console.Green()
    fmt.Println("That appears to have worked, the ami will disspear in a few mins")
    console.Reset()
}
