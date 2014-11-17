package exploud

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	jsonq "github.com/jmoiron/jsonq"
	common "nokia.com/klink/common"
	console "nokia.com/klink/console"
	ditto "nokia.com/klink/ditto"
	onix "nokia.com/klink/onix"
	props "nokia.com/klink/props"
)

const latestVersionString = "latest"

func Init() {
	common.Register(
		common.Component{"deploy", Exploud,
			"{app} {env} [{ami}] Deploy the AMI {ami} for {app} to {env}. (If no ami is specified, the latest is assumed.)", "APPS:ENVS"},
		common.Component{"watch", Watch,
			"{id} Resume watching the deployment with the supplied id", ""},
		common.Component{"undo", Undo,
			"{app} {env} Undo the steps of a broken deployment", "APPS:ENVS"},
		common.Component{"deployments", Deployments,
			"[{app} {env}] Display a list of ongoing deployments or recent deployments if an app is passed", "APPS:ENVS"},
		common.Component{"pause", Pause,
			"{app} {env} attempts to pause a running deployment of {app} in {env}", "APPS:ENVS"},
		common.Component{"cancel-pause", CancelPause,
			"{app} {env} cancels any existing pause for {app} in {env}", "APPS:ENVS"},
		common.Component{"resume", Resume,
			"{app} {env} attempts to resume a paused deployment of {app} in {env}", "APPS:ENVS"},
		common.Component{"rollback", Rollback,
			"{app} {env} rolls the application back to the last successful deploy", "APPS:ENVS"},
		common.Component{"create-app", CreateApp,
			"{app} -E {email} Creates a new application", "APPS:ENVS"},
		common.Component{"boxes", Boxes,
			"{app} {env} -f format [text|json] -S status [stopped|running|terminated]", "APPS:ENVS"})
}

// Returns explouds url with the supplied string appended
func exploudUrl(end string) string {
	return "http://exploud.brislabs.com/1.x" + end
}

// Returns a jsonq object with information about the boxes running
// for the supplied application and environment
func JsonBoxes(app string, env string, i []interface{}) {
	describeUrl := exploudUrl("/describe-instances/" + app + "/" + env)
	common.GetJson(describeUrl, &i)
}

// Return information about the servers running in the supplied environment
func Boxes(args common.Command) {
	if args.SecondPos == "" {
		console.Fail("You must supply an app as the second positional argument")
	}
	app := args.SecondPos
	if args.ThirdPos == "" {
		console.Fail("You must supply an environment as the third positional argument")
	}
	env := args.ThirdPos

	describeUrl := exploudUrl("/describe-instances/" + app + "/" + env)
	if args.Status != "" {
		describeUrl += "?state=" + args.Status
	}

	fmt.Println(common.GetString(describeUrl, func(req *http.Request) {
		if args.Format == "" || args.Format == "text" {
			req.Header.Add("accept", "text/plain")
		}
	}))
}

// AppExists returns true if the application exists according to exploud
func AppExists(appName string) bool {
	return common.Head(exploudUrl("/applications/" + appName))
}

// ************************************************
// **                                            **
// ** Deployment based code, undo, rollback etc  **
// **                                            **
// ************************************************

type AmiDeployRequest struct {
	Ami      string `json:"ami"`
	Message  string `json:"message"`
	Silent   bool   `json:"silent"`
	Username string `json:"user"`
}

type DeployRequest struct {
	Message  string `json:"message"`
	Silent   bool   `json:"silent"`
	Username string `json:"user"`
}

type CreateAppRequest struct {
	Email string `json:"email"`
}

type DeploymentReference struct {
	Id string `json:"id"`
}

// Validate common deployment arguments
func validateDeploymentArgs(args common.Command) {
	app := args.SecondPos
	if app == "" {
		console.Fail("Must supply an application name as the second positional argument.")
	}
	if !AppExists(app) {
		console.Fail(
			fmt.Sprintf("Application \"%s\" does not exist. It's your word against exploud.", app))
	}

	env := args.ThirdPos
	if env == "" {
		console.Fail("Must supply an environment as third postional argument.")
	} else if !onix.KnownEnvironment(env) {
		console.Fail(
			fmt.Sprintf("Third argument \"%s\" must be a known environment. %s.",
				env, onix.GetEnvironments(env)))
	}
}

// Validate common deployment arguments including a message
func validateDeploymentArgsWithMessage(args common.Command) {
	validateDeploymentArgs(args)

	if args.Message == "" {
		console.Fail("Must supply a deploy message using -m")
	}
}

// Validate all deployment arguments including ami as a forth pos argument
func validateDeploymentArgsWithAmi(args common.Command) {
	validateDeploymentArgsWithMessage(args)

	ami := args.FourthPos
	common.FailIfNotAmi(ami)
}

// Execute a deployment
func DoDeployment(url string, body interface{}, args common.Command) {
	deployRef := DeploymentReference{}

	common.PostJsonUnmarshalResponse(url, &body, &deployRef)

	PollDeployNew(deployRef.Id, args.SecondPos)
}

// List any currently active deployments or if an app param supplied
// all deployments for that app
func Deployments(args common.Command) {
	app := args.SecondPos
	env := args.ThirdPos
	if app == "" {
		fmt.Println(common.GetString(exploudUrl("/in-progress")))
	} else {
		if env == "" {
			env = "poke"
		}
		url := exploudUrl("/deployments?application=" + app + "&environment=" + env)
		fmt.Println(common.GetString(url))
	}
}

// Resume an existing deployment, also pretty swish for testing
func Watch(args common.Command) {
	id := args.SecondPos

	if id == "" {
		console.Fail("Must supply a deployment id as the second positional arg")
	}

	PollDeployNew(id, "TODO: do this for a name instead of id")
}

// Exploud -> Expload the app to the cloud. AKA deploy the app named in the args SecondPos
func Exploud(args common.Command) {
	validateDeploymentArgsWithAmi(args)

	app := args.SecondPos
	env := args.ThirdPos
	ami := args.FourthPos

	latestAmi := ditto.LatestAmiFor(app)

	if ami == "" {
		console.Confirmer(
			console.Green,
			fmt.Sprintf("The latest ami %s (version %s) will be deployed. Are you sure you wish to continue?", latestAmi.ImageId, latestAmi.Version))

		ami = latestAmi.ImageId
	} else if latestAmi.ImageId != ami {
		console.Confirmer(
			console.Red,
			fmt.Sprintf("The latest ami for this application is %s (version %s). Are you sure you wish to continue?", latestAmi.ImageId, latestAmi.Version))
	}

	deployUrl := fmt.Sprintf(exploudUrl("/applications/%s/%s/deploy"), app, env)
	deployRequest := AmiDeployRequest{ami, args.Message, args.Silent, props.Get("Username")}

	DoDeployment(deployUrl, deployRequest, args)
}

// Pause a running deployment
func Pause(args common.Command) {
	validateDeploymentArgs(args)

	app := args.SecondPos
	env := args.ThirdPos

	pauseUrl := fmt.Sprintf(exploudUrl("/applications/%s/%s/pause"), app, env)

	fmt.Printf("Attempting to pause deployment of %s in %s\n", app, env)

	common.PostJson(pauseUrl, "")
}

// Cancel a submitted pause
func CancelPause(args common.Command) {
	validateDeploymentArgs(args)

	app := args.SecondPos
	env := args.ThirdPos

	pauseUrl := fmt.Sprintf(exploudUrl("/applications/%s/%s/pause"), app, env)

	fmt.Printf("Attempting to cancel pause of %s in %s\n", app, env)

	common.Delete(pauseUrl)
}

// Resume a paused deployment
func Resume(args common.Command) {
	validateDeploymentArgs(args)

	app := args.SecondPos
	env := args.ThirdPos

	resumeUrl := fmt.Sprintf(exploudUrl("/applications/%s/%s/resume"), app, env)

	fmt.Printf("Attempting to resume deployment of %s in %s\n", app, env)

	common.PostJson(resumeUrl, "")
}

// Undo the steps from a borked deployment
func Undo(args common.Command) {
	validateDeploymentArgsWithMessage(args)

	app := args.SecondPos
	env := args.ThirdPos

	deployUrl := fmt.Sprintf(exploudUrl("/applications/%s/%s/undo"), app, env)
	deployRequest := DeployRequest{args.Message, args.Silent, props.Get("Username")}

	DoDeployment(deployUrl, deployRequest, args)
}

// Exploud -> Expload the app to the cloud. AKA deploy the app named in the args SecondPos
// Must pass SecondPos and Ami arguments
func Rollback(args common.Command) {
	validateDeploymentArgsWithMessage(args)

	app := args.SecondPos
	env := args.ThirdPos

	deployUrl := fmt.Sprintf(exploudUrl("/applications/%s/%s/rollback"), app, env)
	deployRequest := DeployRequest{args.Message, args.Silent, props.Get("Username")}

	DoDeployment(deployUrl, deployRequest, args)
}

// Returns the status of the deployment with the supplied id
func GetDeploymentStatus(deploymentId string, retries int) string {
	url := exploudUrl(fmt.Sprintf("/deployments/%s", deploymentId))

	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error creating GET request for url: %s\n", url)
		panic(err)
	}
	if resp.StatusCode == 404 {
		if retries > 0 {
			time.Sleep(1 * time.Second)
			return GetDeploymentStatus(deploymentId, retries-1)
		}
	}

	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	data := map[string]interface{}{}
	dec := json.NewDecoder(strings.NewReader(string(responseBody)))
	dec.Decode(&data)
	responsejq := jsonq.NewQuery(data)

	status, err := responsejq.String("status")
	if err != nil {
		fmt.Println("Failed to parse deployment status")
		panic(err)
	}
	return status
}

// Prints out the status line for the deploy
func PrintStatus(taskId string, serviceName string, status string) {
	// first line
	fmt.Println("")
	console.Red()
	console.Bold()
	fmt.Print(fmt.Sprintf("%30s", "Explouding "+serviceName+": "))
	console.Green()
	fmt.Println(taskId)

	// status line
	console.Red()
	fmt.Print(fmt.Sprintf("%30s", "Status: "))
	if status == "completed" {
		console.Green()
	} else {
		console.Brown()
	}
	fmt.Println(status)
	console.FReset()
}

// Prints any new logs since the previous lastTime
// if lastTime is blank then return all logs
// Returns the new lastTime
func PrintNewDeploymentLogs(deploymentId string, lastTime string) string {
	url := exploudUrl(fmt.Sprintf("/deployments/%s/logs", deploymentId))
	if lastTime != "" {
		url += "?since=" + lastTime
	}

	logs, err := common.GetAsJsonq(url).ArrayOfObjects("logs")
	if err != nil {
		fmt.Println("Could not parse deployment logs from: " + url)
		panic(err)
	}
	for _, log := range logs {
		logjq := jsonq.NewQuery(log)
		// naughty... TODO: put in some err handling
		lastTime, _ = logjq.String("date")
		message, _ := logjq.String("message")
		fmt.Println(message)
	}
	return lastTime
}

// Poll the supplied deployment printing out the status to the console.
func PollDeployNew(deploymentId string, serviceName string) {
	chnl := HandleDeployInterrupt()
	defer DeregisterInterupt(chnl)

	PrintStatus(deploymentId, serviceName, "pending")

	status := GetDeploymentStatus(deploymentId, 10)

	lastTime := ""
	for true {
		lastTime = PrintNewDeploymentLogs(deploymentId, lastTime)

		if status == "failed" || status == "invalid" {
			console.Red()
			console.Fail("Deployment reached a failed or terminated status :-(")
		}

		if status == "completed" {
			break
		}

		// continue
		time.Sleep(5 * time.Second)
		status = GetDeploymentStatus(deploymentId, 1)
	}

	console.Green()
	PrintStatus(deploymentId, serviceName, "Finished!")
	console.Reset()
}

// Register a new application with exploud, should have the knock on effect
// of registering with the other applications that exploud depends upon e.g.
// onix and tyranitar
func CreateApp(args common.Command) {
	if args.SecondPos == "" || strings.Index(args.SecondPos, "-") == 0 {
		console.Fail("Must supply an application name as second positional argument")
	}

	if args.Email == "" {
		console.Fail("Don't be lazy! You must supply a value for email")
	}

	fmt.Printf(
		"Calling exploud to create application %s with email: %s",
		args.SecondPos,
		args.Email,
	)

	createBody := CreateAppRequest{args.Email}

	response := common.PutJson(exploudUrl("/applications/"+args.SecondPos), createBody)

	fmt.Println("Exploud has created our application for us!")
	fmt.Println(response)
}

// Interrupt constants
const (
	Yes = iota
	No
	Continue
)

// Returns true if the user wants to cancel the deployment
func cancelDeploymentPerchance() int {
	console.Red()
	fmt.Println("CURRENTLY THIS DOESN'T WORK. COMING SOON 2009!")
	fmt.Println("If you kill the deployment here it will continue")
	fmt.Println("You can start watching it again with klink watch {app} {env}")
	fmt.Println("Do you want to rollback the deployment? [Yes, No, Continue]")
	fmt.Println("CURRENTLY THIS DOESN'T WORK. COMING SOON 2009!")
	console.Reset()
	var response string

	fmt.Scan(&response)

	switch response {
	case "yes", "Yes", "YES", "y", "Y":
		return Yes
	case "no", "No", "NO", "n", "N":
		return No
	case "continue", "cont", "Continue", "c", "C":
		return Continue
	default:
		fmt.Println("Type better.")
		return cancelDeploymentPerchance()
	}
}

// Handle interupts and ask the user if they want to rollback the deployment
func HandleDeployInterrupt() chan os.Signal {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			fmt.Println(sig)
			switch cancelDeploymentPerchance() {
			case Yes:
				fmt.Println("This will rollback the deployment when exploud is ready!")
				os.Exit(0)
			case No:
				fmt.Println("Not rollingback. Just exiting. Your deployment will continue.")
				os.Exit(1)
			case Continue:
				fmt.Println("Continuing...")
			}
		}
	}()
	return c
}

// Deregister the interupt
func DeregisterInterupt(c chan<- os.Signal) {
	signal.Stop(c)
}
