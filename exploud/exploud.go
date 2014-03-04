package exploud

import (
	"fmt"
	"net/http"
	common "nokia.com/klink/common"
	console "nokia.com/klink/console"
	props "nokia.com/klink/props"
	ditto "nokia.com/klink/ditto"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"time"
)

const latestVersionString = "latest"

func Init() {
	common.Register(
		common.Component{"deploy", Exploud,
			"{app} {env} [{ami}] Deploy the AMI {ami} for {app} to {env}. (If no ami is specified, the latest is assumed.)"},
		common.Component{"undo", Undo,
			"{app} {env} Undo the steps of a broken deployment"},
		common.Component{"rollback", Rollback,
			"{app} {env} rolls the application back to the last successful deploy"},
		common.Component{"list-apps", ListApps,
			"Lists the applications that exist (via exploud)"},
		common.Component{"create-app", CreateApp,
			"{app} -E {email} -o {owner} -d {description} Creates a new application"},
		common.Component{"boxes", Boxes,
			"{app} {env} -f format [text|json] -S status [stopped|running|terminated]"})
}

// Returns explouds url with the supplied string appended
func exploudUrl(end string) string {
	return "http://exploud.brislabs.com:8080/1.x" + end
}

// Return information about the servers running in the supplied environment
func Boxes(args common.Command) {
	if args.SecondPos == "" {
		console.Fail("You must supply a service as the second positional argument")
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
		if args.Format == "text" {
			req.Header.Add("accept", "text/plain")
		}
	}))
}

// List the apps known by exploud
func ListApps(args common.Command) {
	fmt.Println(common.GetString(exploudUrl("/applications")))
}

// AppExists returns true if the application exists according to the exploud service
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
	Username string `json:"user"`
}

type DeployRequest struct {
	Message  string `json:"message"`
	Username string `json:"user"`
}

type CreateAppRequest struct {
	Description string `json:"description"`
	Email       string `json:"email"`
	Owner       string `json:"owner"`
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
	} else if !(env == "poke" || env == "prod") {
		console.Fail(
			fmt.Sprintf("Third argument \"%s\" must be an environment. poke or prod.", env))
	}

	if args.Message == "" {
		console.Fail("Must supply a deploy message using -m")
	}
}

// Validate all deployment arguments including ami as a forth pos argument
func validateDeploymentArgsWithAmi(args common.Command) {
	validateDeploymentArgs(args)

	ami := args.FourthPos
	if ami != "" {
		matched, err := regexp.MatchString("^ami-.+$", ami)
		if err != nil {
			panic(err)
		}
		if !matched {
			console.Fail(fmt.Sprintf("%s Doesn't look like an ami", ami))
		}
	}
}

// Execute a deployment
func DoDeployment(url string, body interface{}, message string, args common.Command) {
	deployRef := DeploymentReference{}

	common.PostJsonUnmarshalResponse(url, &body, &deployRef)

	console.Hubot(message, args)

	PollDeployNew(deployRef.Id, args.SecondPos)
}

// Exploud -> Expload the app to the cloud. AKA deploy the app named in the args SecondPos
func Exploud(args common.Command) {
	validateDeploymentArgsWithAmi(args)

	app := args.SecondPos
	env := args.ThirdPos
	ami := args.FourthPos

	latestAmi := ditto.LatestAmiFor(app)

	if (ami == "") {
		confirmDeployLatest(latestAmi)
		ami = latestAmi.ImageId
	} else if (latestAmi.ImageId != ami) {
		confirmNonLatestBake(latestAmi)
	}

	deployUrl := fmt.Sprintf(exploudUrl("/applications/%s/%s/deploy"), app, env)
	deployRequest := AmiDeployRequest{ami, args.Message, props.GetUsername()}
	message := fmt.Sprintf("%s is deploying %s for service %s to %s. %s",
		props.GetUsername(), ami, app, env, args.Message)

	DoDeployment(deployUrl, deployRequest, message, args)
}

func confirmNonLatestBake(ami ditto.Ami) {

	console.Red()
	fmt.Println(fmt.Sprintf("The latest ami for this application is %s (version %s). Are you sure you wish to continue?", ami.ImageId, ami.Version))
	console.Reset()

	var response string

	fmt.Scan(&response)

	switch response {
	case "yes", "Yes", "YES", "y", "Y":
		break
	case "no", "No", "NO", "n", "N":
		console.Red()
		console.Fail("Deployment aborted.")
		console.Reset()
	default:
		fmt.Println("Type better.")
		confirmNonLatestBake(ami)
	}
}

func confirmDeployLatest(latestAmi ditto.Ami) {

	console.Green()
	fmt.Println(fmt.Sprintf("The latest ami %s (version %s) will be deployed. Are you sure you wish to continue?", latestAmi.ImageId, latestAmi.Version))
	console.Reset()

	var response string

	fmt.Scan(&response)

	switch response {
	case "yes", "Yes", "YES", "y", "Y":
		break
	case "no", "No", "NO", "n", "N":
		console.Red()
		console.Fail("Deployment aborted.")
		console.Reset()
	default:
		fmt.Println("Type better.")
		confirmDeployLatest(latestAmi)
	}
}

// Undo the steps from a borked deployment
func Undo(args common.Command) {
	validateDeploymentArgs(args)

	app := args.SecondPos
	env := args.ThirdPos

	deployUrl := fmt.Sprintf(exploudUrl("/applications/%s/%s/undo"), app, env)
	deployRequest := DeployRequest{args.Message, props.GetUsername()}
	message := fmt.Sprintf("%s is undoing deployment of service %s in %s. %s",
		props.GetUsername(), app, env, args.Message)

	DoDeployment(deployUrl, deployRequest, message, args)
}

// Exploud -> Expload the app to the cloud. AKA deploy the app named in the args SecondPos
// Must pass SecondPos and Ami arguments
func Rollback(args common.Command) {
	validateDeploymentArgs(args)

	app := args.SecondPos
	env := args.ThirdPos

	deployUrl := fmt.Sprintf(exploudUrl("/applications/%s/%s/rollback"), app, env)
	deployRequest := DeployRequest{args.Message, props.GetUsername()}
	message := fmt.Sprintf("%s is rollingback service %s in %s. %s",
		props.GetUsername(), app, env, args.Message)

	DoDeployment(deployUrl, deployRequest, message, args)
}

// Exploud JSON task log message
type TaskLog struct {
	Date    string `json:"Date"`
	Message string `json:"message"`
}

// Exploud JSON task
type Task struct {
	Action         string    `json:"action"`
	DurationString string    `json:"durationString"`
	End            string    `json:"end"`
	Id             string    `json:"_id"`
	Log            []TaskLog `json:"log"`
	Operation      string    `json:"operation"`
	Start          string    `json:"start"`
	Status         string    `json:"status"`
	Url            string    `json:"url"`
}

// Exploud JSON deployment
type Deployment struct {
	Ami         string `json:"ami"`
	Application string `json:"application"`
	Created     string `json:"created"`
	End         string `json:"end"`
	Environment string `json:"environment"`
	Hash        string `json:"hash"`
	Id          string `json:"id"`
	Region      string `json:"region"`
	Start       string `json:"start"`
	Tasks       []Task `json:"tasks"`
	User        string `json:"user"`
}

func GetDeployment(deploymentId string) Deployment {
	url := exploudUrl(fmt.Sprintf("/deployments/%s", deploymentId))

	deployment := Deployment{}
	common.GetJson(url, &deployment)
	return deployment
}

// Prints out the status line for the deploy
func Status(taskId string, serviceName string, status string) {
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

// Poll the supplied deployment printing out the status to the console.
func PollDeployNew(deploymentId string, serviceName string) {
	chnl := HandleDeployInterrupt()
	defer DeregisterInterupt(chnl)

	Status(deploymentId, serviceName, "pending")
	deployment := GetDeployment(deploymentId)

	timeout := time.Now().Add((20 * time.Minute))
	for i := 0; i < len(deployment.Tasks) && time.Now().Before(timeout); i++ {
		task := deployment.Tasks[i]

		console.Green()
		fmt.Printf("Starting task: %s\n\n", task.Action)
		console.Reset()

		previousLength := 0
		// can't check == running as wont be set when we first call
		for (task.Status != "completed") &&
			(task.Status != "undone") &&
			(task.Status != "failed") &&
			(task.Status != "skipped") &&
			(task.Status != "teminated") &&
			time.Now().Before(timeout) {

			time.Sleep(5 * time.Second)
			deployment = GetDeployment(deploymentId)
			task = deployment.Tasks[i]

			for i := previousLength; i < len(task.Log); i++ {
				fmt.Println(task.Log[i])
			}

			previousLength = len(task.Log)
		}

		// if we see something failed then kill everything - exploud doesn't recover
		if task.Status == "failed" || task.Status == "terminated" {
			console.Fail(fmt.Sprintf("Deployment reached a failed or terminated task: %s", task))
		}
	}

	console.Green()
	Status(deploymentId, serviceName, "Finished!")
	console.Reset()
}

// Register a new application with exploud, should have the knock on effect
// of registering with the other services that exploud depends upon e.g.
// onix and tyranitar
func CreateApp(args common.Command) {
	if args.SecondPos == "" || strings.Index(args.SecondPos, "-") == 0 {
		console.Fail("Must supply an application name as second positional argument")
	}

	if args.Description == "" || args.Email == "" || args.Owner == "" {
		console.Fail("Don't be lazy! You must supply owner, email and description values")
	}

	fmt.Printf(
		"Calling exploud to create application %s with description:\n%s, email: %s, owner: %s",
		args.SecondPos,
		args.Description,
		args.Email,
		args.Owner,
	)

	createBody := CreateAppRequest{args.Description, args.Email, args.Owner}

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
	fmt.Println("Do you want to rollback the deployment? [Yes, No, Continue]")
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
