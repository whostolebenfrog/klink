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
		console.Fail("Must supply an application name as the second positional argument")
	}
	if args.Ami == "" {
		console.Fail("Must supply an ami to deploy using --ami")
	}
	if !AppExists(args.SecondPos) {
		console.Fail(fmt.Sprintf("Application \"%s\" does not exist. It's your word aginst exploud!",
			args.SecondPos))
	}

	deployUrl := fmt.Sprintf(exploudUrl("/applications/%s/deploy"), args.SecondPos)

	deployRequest := DeployRequest{args.Ami, "dev"}
	task := TaskReference{}

	err := common.PostJsonUnmarshalResponse(deployUrl, &deployRequest, &task)
	if err != nil {
		fmt.Println(err)
		console.Fail("Error calling exploud, exiting.")
	}

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

	err := common.GetJson(taskUrl, &task)
	if err != nil {
		fmt.Println(err)
		console.Fail("Unable to get task from exploud")
	}
	return task
}

// TODO - lining up depends on version length! BENKINS!
// Prints out the status line for the deploy
func Status(taskId string, serviceName string, status string) {
	fmt.Print(fmt.Sprintf("\033[1m\033[31m    Explouding %s: ", serviceName))
	fmt.Println("\033[32m", taskId)
	statusColor := 33
	if status == "completed" {
		statusColor = 32
	}
	fmt.Println(fmt.Sprintf("\033[1m\033[31m              Status:\033[%dm  %s",
		statusColor, status))
	fmt.Println("\033[0m")
}

// Poll the supplied taskId, printing the status to the console. Finishing
// after either the task is marked as completed or a timeout is reached
func PollDeploy(taskId string, serviceName string) {
	fmt.Println("\n")
	Status(taskId, serviceName, "pending")
	task := GetTask(taskId)

	timeout := time.Now().Add((20 * time.Minute))
	previousLength := 0
	for (task.Status != "completed") &&
		(task.Status != "failed") &&
		(task.Status != "teminated") &&
		time.Now().Before(timeout) {

		time.Sleep(5 * time.Second)
		task = GetTask(taskId)

		// Jump the cursor up the right number of lines and clear
		for i := 0; i < previousLength+4; i++ {
			fmt.Println("\033[2A\033[2K\r")
		}
		Status(taskId, serviceName, task.Status)

		previousLength = len(task.Log) - 1
		for line := range task.Log {
			fmt.Println(task.Log[line])
		}
	}

	fmt.Print("\033[0m")
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

	response, err := common.PutJson(exploudUrl("/applications/"+args.SecondPos), createBody)

	if err != nil {
		fmt.Println(err, response)
		console.BigFail("Unable to register new application with exploud")
	}

	fmt.Println("Exploud has created our application for us!")
	fmt.Println(response)
}

// List the apps known by exploud
func ListApps() {
	response, err := common.GetString(exploudUrl("/applications"))
	if err != nil {
		fmt.Println(response, err)
		console.Fail("Error listing applications")
	}
	fmt.Println(response)
}

// AppExists returns true if the application exists according to the exploud service
func AppExists(appName string) bool {
	resp, _ := common.Head(exploudUrl("/applications/" + appName))
	return resp
}
