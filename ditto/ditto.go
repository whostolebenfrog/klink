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
	jenkins "nokia.com/klink/jenkins"
	onix "nokia.com/klink/onix"
)

func Init() {
	// TODO  --- helpers file! the first one
	common.Register(
		common.Component{"ditto", Helpers,
			"Various helpers; lock, unlock, clean, build public and ent base amis", "DITTOS"},
		common.Component{"bake", LiveBake,
			"{app} {version} [-t {hvm,para}] Bakes an AMI for {app} with version {version}", "APPS"},
		common.Component{"betabake", BetaBake,
			"{app} {version} [-t {hvm,para}] Bakes an Amazon Linux AMI for {app} with version {version}", "APPS"},
		common.Component{"allow-prod", AllowProd,
			"{app} Allows the prod aws account access to the supplied application", "APPS"},
		common.Component{"amis", FindAmis,
			"{app} Lists the latest amis for the supplied application name", "APPS"},
		common.Component{"images", FindAmis,
			"{app} Lists the latest images for the supplied application name", "APPS"},
		common.Component{"latest-bake", LatestBake,
			"{app} Outputs the latest baked version of the specified application", "APPS"},
		common.Component{"delete-ami", DeleteAmi,
			"{app} {ami} Removes the supplied ami, makes it undeployable.", "APPS"})
}

func dittoUrl(end string) string {
	base := os.Getenv("DITTO_URL")
	if base == "" {
		base = "http://ditto.brislabs.com:8080"
	}
	return base + "/1.x" + end
}

func betaDittoUrl(end string) string {
	return "http://internal-betaditto-2028158683.eu-west-1.elb.amazonaws.com:8080/1.x" + end
}

func bakeUrl(app string, version string) string {
	return fmt.Sprintf(dittoUrl("/bake/%s/%s"), app, version)
}

func betaBakeUrl(app string, version string) string {
	return fmt.Sprintf(betaDittoUrl("/bake/%s/%s"), app, version)
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

type bakeUrlFn func(string, string) string

// Do the bake with the supplied url lookup function
func Bake(args common.Command, bUrl bakeUrlFn) {
	app := args.SecondPos
	version := args.ThirdPos
	virtType := args.Type

	if app == "" {
		console.Fail("Application must be supplied as the second argument")
	}

	if !onix.AppExists(app) {
		console.Fail(fmt.Sprintf("Application '%s' does not exist. It's your word against onix!",
			app))
	}

	if version == "" {
		latestVersion, dateDescription := jenkins.GetLatestStableBuildVersion(jenkins.JobPath(app, "releasePath"))
		if latestVersion != "" {
			console.Confirmer(
				console.Green,
				fmt.Sprintf("Version %s built on %s will be baked. Are you sure you wish to continue?", latestVersion, dateDescription))

			version = latestVersion
		}
	}

	if version == "" {
		console.Fail("Version must be supplied as the third argument or using --version")
	}

	url := bUrl(app, version)
	if args.Type != "" {
		url += "?virt-type=" + virtType
	} else {
		bakeType := onix.GetOptionalProperty(app, "bakeType")
		if bakeType != "" {
			url += "?virt-type=" + bakeType
		}
	}
	DoBake(url, 120)
}

// BetaBake the ami - e.g. use the version of ditto in beta.
func BetaBake(args common.Command) {
	Bake(args, betaBakeUrl)
}

// Bake the ami
func LiveBake(args common.Command) {
	Bake(args, bakeUrl)
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

// list of helper names for autocomplete, could tie these to functions and combine with
// this list below but might just make it more complex for little need.
func HelperNames() []string {
	return []string{"lock", "unlock", "clean", "entertainment", "public", "inprogress"}
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
		unlockUrl := dittoUrl("/lock")
		common.Delete(unlockUrl)
		console.Green()
		fmt.Println("unlocked")
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
	fmt.Println("That appears to have worked, the ami will disappear in a few minutes")
	console.Reset()
}
