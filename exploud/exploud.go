package exploud

import (
	"fmt"
	common "nokia.com/klink/common"
	console "nokia.com/klink/console"
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

type TaskReference struct {
	TaskId string `json:"taskId"`
}

// Exploud -> Expload the app to the cloud. AKA deploy the app named in the args SecondPos
// Must pass SecondPos and Ami arguments
func Exploud(args common.Command) {
	if args.SecondPos == "" {
		console.Fail("Must supply an application name as the second positional argument.")
	}
	if args.ThirdPos == "" {
		console.Fail("Must supply an environment as third postional argument. Use dev or prod.")
	}
	if args.ForthPos == "" {
		if args.Ami == "" {
			console.Fail("Must supply an ami as the forth argument or with --ami.")
		}
		args.ForthPos = args.Ami
	}
	if !AppExists(args.SecondPos) {
		console.Fail(fmt.Sprintf("Application \"%s\" does not exist. It's your word aginst exploud.",
			args.SecondPos))
	}

	deployUrl := fmt.Sprintf(exploudUrl("/applications/%s/deploy"), args.SecondPos)

	deployRequest := DeployRequest{args.ForthPos, args.ThirdPos}
	task := TaskReference{}

	common.PostJsonUnmarshalResponse(deployUrl, &deployRequest, &task)

	// TODO: user and version (can be parsed from the ami name)
	hubotMessage := fmt.Sprintf("Deploying %s for service %s to %s.",
		args.ForthPos, args.SecondPos, args.ThirdPos)
	if args.Message != "" {
		hubotMessage += " " + args.Message + "."
	}
	console.Hubot(hubotMessage, args)

	PollDeploy(task.TaskId, args.SecondPos)
}

// Exploud JSON task log message
type TaskLog struct {
	Date    string `json:"Date"`
	Message string `json:"message"`
}

// Exploud JSON task
type Task struct {
	DurationString string    `json:"durationString"`
	Id             string    `json:"_id"`
	Log            []TaskLog `json:"log"`
	Operation      string    `json:"operation"`
	Region         string    `json:"region"`
	RunId          string    `json:"runId"`
	Status         string    `json:"status"`
	UpdateTime     string    `json:"updateTime"`
	WorkflowId     string    `json:"workflowId"`
}

// Get the task for the supplied id
func GetTask(taskId string) Task {
	taskUrl := exploudUrl(fmt.Sprintf("/tasks/%s", taskId))

	task := Task{}
	common.GetJson(taskUrl, &task)
	return task
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

// Poll the supplied taskId, printing the status to the console. Finishing
// after either the task is marked as completed or a timeout is reached
func PollDeploy(taskId string, serviceName string) {
	Status(taskId, serviceName, "pending")
	task := GetTask(taskId)

	timeout := time.Now().Add((20 * time.Minute))
	previousLength := 0
	// can't check == running as wont be set when we first call
	for (task.Status != "completed") &&
		(task.Status != "failed") &&
		(task.Status != "teminated") &&
		time.Now().Before(timeout) {

		time.Sleep(5 * time.Second)
		task = GetTask(taskId)

		for i := previousLength; i < len(task.Log); i++ {
			fmt.Println(task.Log[i])
		}

		previousLength = len(task.Log)
	}

	Status(taskId, serviceName, task.Status)
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
