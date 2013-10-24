package exploud

import (
	"fmt"
	common "nokia.com/klink/common"
	console "nokia.com/klink/console"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"time"
)

type DeployRequest struct {
	Ami         string `json:"ami"`
	Environment string `json:"environment"`
}

type CreateAppRequest struct {
	Description string `json:"description"`
	Email       string `json:"email"`
	Owner       string `json:"owner"`
}

func exploudUrl(end string) string {
	return "http://exploud.brislabs.com:8080/1.x" + end
}

type DeploymentReference struct {
	Id string `json:"id"`
}

func validateDeploymentArgs(args common.Command) {
	if args.SecondPos == "" {
		console.Fail("Must supply an application name as the second positional argument.")
	}
	if !AppExists(args.SecondPos) {
		console.Fail(fmt.Sprintf("Application \"%s\" does not exist. It's your word aginst exploud.",
			args.SecondPos))
	}

	if args.ThirdPos == "" {
		console.Fail("Must supply an environment as third postional argument.")
	} else if !(args.ThirdPos == "poke" || args.ThirdPos == "prod") {
		console.Fail(fmt.Sprintf("Third argument \"%s\" must be an environment. poke or prod.",
			args.ThirdPos))
	}

	if args.FourthPos == "" {
		if args.Ami == "" {
			console.Fail("Must supply an ami as the fourth argument or with --ami.")
		}
		args.FourthPos = args.Ami
	}
	matched, err := regexp.MatchString("^ami-.+$", args.FourthPos)
	if err != nil {
		panic(err)
	}
	if !matched {
		console.Fail(fmt.Sprintf("%s Doesn't look like an ami", args.FourthPos))
	}
}

// Exploud -> Expload the app to the cloud. AKA deploy the app named in the args SecondPos
// Must pass SecondPos and Ami arguments
func Exploud(args common.Command) {
	validateDeploymentArgs(args)

	deployUrl := fmt.Sprintf(exploudUrl("/applications/%s/deploy"), args.SecondPos)

	deployRequest := DeployRequest{args.FourthPos, args.ThirdPos}
	deployRef := DeploymentReference{}

	common.PostJsonUnmarshalResponse(deployUrl, &deployRequest, &deployRef)

	// TODO: user and version (can be parsed from the ami name)
	hubotMessage := fmt.Sprintf("Deploying %s for service %s to %s.",
		args.FourthPos, args.SecondPos, args.ThirdPos)
	if args.Message != "" {
		hubotMessage += " " + args.Message + "."
	}

    // TODO: put this back
	//console.Hubot(hubotMessage, args)
	PollDeployNew(deployRef.Id, args.SecondPos)
}

// Exploud JSON task log message
type TaskLog struct {
	Date    string `json:"Date"`
	Message string `json:"message"`
}

// Exploud JSON task
type Task struct {
    Action string `json:"action"`
	DurationString   string        `json:"durationString"`
	End              string        `json:"end"`
	Id               string        `json:"_id"`
	Log              []TaskLog     `json:"log"`
	Operation        string        `json:"operation"`
	Start            string        `json:"start"`
	Status           string        `json:"status"`
	Url              string        `json:"url"`
}

// Exploud JSON deployment
type Deployment struct {
	Ami         string        `json:"ami"`
	Application string        `json:"application"`
	Created     string        `json:"created"`
	End         string        `json:"end"`
	Environment string        `json:"environment"`
	Hash        string        `json:"hash"`
	Id          string        `json:"id"`
	Region      string        `json:"region"`
	Start       string        `json:"start"`
	Tasks       []Task        `json:"tasks"`
	User        string        `json:"user"`
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
        fmt.Println(fmt.Sprintf("Starting task: %s\n", task.Action))
        console.Reset()

        previousLength := 0
        // can't check == running as wont be set when we first call
        for (task.Status != "completed") &&
            (task.Status != "failed") &&
            (task.Status != "teminated") &&
            time.Now().Before(timeout) {

            // if we see something failed then kill everything - exploud doesn't recover
            if task.Status == "failed" || task.Status == "terminated" {
                console.Fail(fmt.Sprintf("Deployment reached a failed or terminated task: %s", task))
            }

            time.Sleep(5 * time.Second)
            deployment = GetDeployment(deploymentId)
            task = deployment.Tasks[i]

            for i := previousLength; i < len(task.Log); i++ {
                fmt.Println(task.Log[i])
            }

            previousLength = len(task.Log)
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

	fmt.Println(fmt.Sprintf(`Calling exploud to create application %s with description:
        %s, email: %s, owner: %s`,
		args.SecondPos, args.Description, args.Email, args.Owner))

	createBody := CreateAppRequest{args.Description, args.Email, args.Owner}

	response := common.PutJson(exploudUrl("/applications/"+args.SecondPos), createBody)

	fmt.Println("Exploud has created our application for us!")
	fmt.Println(response)
}

// List the apps known by exploud
func ListApps() {
	fmt.Println(common.GetString(exploudUrl("/applications")))
}

// AppExists returns true if the application exists according to the exploud service
func AppExists(appName string) bool {
	return common.Head(exploudUrl("/applications/" + appName))
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
