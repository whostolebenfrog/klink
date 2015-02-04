package jenkins

import (
	"fmt"
	"math"
	"net/http"
	"os"
	"strings"
	"time"

	jsonq "github.com/jmoiron/jsonq"
	common "mixrad.io/klink/common"
	console "mixrad.io/klink/console"
	onix "mixrad.io/klink/onix"
)

func Init() {
	common.Register(
		common.Component{"build", Build,
			"{app} runs the Jenkins release job for an application", "APPS"},
		common.Component{"test", Test,
			"{app} runs the Jenkins test job for an application", "APPS"},
		common.Component{"jobs", Jobs, "{app} lists the set of jenkins jobs for an application with their current state", "APPS"})
}

// Build a release job for the supplied application and poll the reponse
func Build(args common.Command) {
	app := args.SecondPos
	if app == "" {
		console.Fail("Yeah, you're gonna have to tell me what to build...")
	}

	path := jobPath(app, "releasePath") + "build"
	fmt.Println("URL: " + path)

	CreateBuild(path)
}

// A unparameterised build has no actions, return true if the build
// is parameterised
func isBuildWithParams(path string) bool {
	actions, err := common.GetAsJsonq(path + "/api/json").Array("actions")
	if err != nil {
		return false
	}
	for _, action := range actions {
		if action.(map[string]interface{})["parameterDefinitions"] != nil {
			return true
		}
	}
	return false
}

// Build a test job for the supplied application and poll the reponse
func Test(args common.Command) {
	app := args.SecondPos
	if app == "" {
		console.Fail("Yeah, you're gonna have to tell me what to test...")
	}

	path := jobPath(app, "testPath")
	if isBuildWithParams(path) {
		path += "buildWithParameters"
	} else {
		path += "build"
	}
	CreateBuild(path)
}

// Create a new build for the given job and poll the reponse
func CreateBuild(path string) {
	location := PostBuild(path) + "api/json"
	job := GetJobFromQueue(location, 12)

	console.Green()
	fmt.Println("\nBuild started, polling Jenkins for output...\n")
	console.Reset()

	status := PollBuild(job)

	if status != "SUCCESS" {
		console.Fail(fmt.Sprintf("Jenkins job failed with status of %s", status))
	}
}

// Returns the release path for the supplied app
func jobPath(app string, property string) string {
	path := onix.GetProperty(app, property)
	if !strings.HasSuffix(path, "/") {
		path += "/"
	}

	return path
}

func GetLatestStableBuildVersion(app string) (string, string) {
	path := jobPath(app, "releasePath") + "api/json"
	builds := common.GetAsJsonq(path)
	buildUrl, err := builds.String("lastStableBuild", "url")

	if err != nil {
		fmt.Println("Unable to parse Jenkins response.")
		panic(err)
	}

	latestStableBuild := common.GetAsJsonq(buildUrl + "/api/json")
	timestamp, err1 := latestStableBuild.Int("timestamp")
	description, err2 := latestStableBuild.String("description")

	if err1 != nil || err2 != nil {
		return "", ""
	} else {
		return description, timestampDescription(timestamp)
	}
}

func timestampDescription(timestamp int) string {
	return time.Unix(0, int64(timestamp)*int64(time.Millisecond)).Format("Mon, Jan 2 15:04:05")
}

// Post a build and return the jobs location in the queue
func PostBuild(url string) string {
	resp, err := http.Post(url, "", nil)
	if err != nil {
		fmt.Printf("Error posting build job to path: %s\n", url)
		panic(err)
	}
	defer resp.Body.Close()

	return resp.Header.Get("Location")
}

// Return the actual job once resolved from the queue, accepts a number of retries
// as the job won't return if it's queued or in the 'quiet period'
func GetJobFromQueue(path string, retries int) string {
	jq := common.GetAsJsonq(path)
	obj, err := jq.String("executable", "url")
	if err != nil {
		if retries > 0 {
			fmt.Println("Jenkins may be queued or in quiet period, retrying in 5 seconds")
			time.Sleep(5 * time.Second)
			return GetJobFromQueue(path, retries-1)
		} else {
			fmt.Println("Unable to parse Jenkins reponse, build may be in a queue")
			panic(err)
		}
	}

	return obj
}

// Poll a build and print the status
func PollBuild(path string) string {
	status := GetJobStatus(path)
	offset := 0
	timeout := time.Now().Add((20 * time.Minute))

	for (status == "in progress...") && time.Now().Before(timeout) {
		status = GetJobStatus(path)
		output, outputBytes := GetJobOutput(path, offset)
		os.Stdout.Write(output)
		offset = outputBytes

		time.Sleep(1 * time.Second)
	}

	fmt.Println(status)
	return status
}

// Return as jobs status
func GetJobStatus(path string) string {
	path += "api/json"
	jq := common.GetAsJsonq(path)
	status, err := jq.String("result")
	if err != nil {
		_, err := jq.String("id")
		if err != nil {
			// case where we definitely don't have a good response
			fmt.Println("Unable to get status from response. Check your build manually.")
			panic(err)
		}
		// returns null for result when in progress which causes an err to be thrown
		// both jenkins and jsonq for go suck
		return "in progress..."
	}
	return status
}

func GetJobOutput(path string, offset int) ([]byte, int) {
	path += fmt.Sprintf("logText/progressiveText?start=%d", offset)
	return common.GetBytes(path)
}

func Jobs(args common.Command) {
	app := args.SecondPos
	if app == "" {
		console.Fail("You didn't supply an app.")
	}

	url := jobPath(app, "jobsPath") + "api/json?depth=2"
	jq := common.GetAsJsonq(url)
	jobs, err := jq.ArrayOfObjects("jobs")

	if err != nil {
		fmt.Println("Couldn't parse the jobs response from: " + url)
		panic(err)
	} else {
		for _, job := range jobs {
			PrintJob(job)
		}
	}
}

func PrintJob(job map[string]interface{}) {
	jobJq := jsonq.NewQuery(job)
	name, _ := jobJq.String("name")
	color, _ := jobJq.String("color")
	lastBuildMs, _ := jobJq.Int("lastBuild", "timestamp")

	if strings.HasSuffix(color, "_anime") {
		console.Bold()
	}

	if strings.HasPrefix(color, "blue") {
		console.Cyan()
		fmt.Print(name)
	} else if strings.HasPrefix(color, "yellow") {
		console.Yellow()
		fmt.Print(name)
	} else if strings.HasPrefix(color, "red") {
		console.Red()
		fmt.Print(name)
	} else {
		console.Grey()
		fmt.Print(name)
	}
	console.Reset()
	fmt.Print(" (last build: " + LastBuildText(lastBuildMs) + ")")
	fmt.Println()
}

func LastBuildText(lastBuildMs int) string {
	if lastBuildMs == 0 {
		return "N/A"
	}

	buildTime := time.Unix(int64(lastBuildMs)/1000, 0)
	now := time.Now()
	diff := now.Sub(buildTime)
	hours := math.Floor(diff.Hours())
	minutes := math.Floor(diff.Minutes())

	if hours >= 48 {
		days := math.Floor(hours / 24)
		hours = hours - (24 * days)
		return fmt.Sprintf("%d days %d hr", int(days), int(hours))
	} else if hours >= 24 {
		return fmt.Sprintf("1 day %d hr", int(hours))
	} else {
		minutes = minutes - (60 * hours)
		return fmt.Sprintf("%d hr %d min", int(hours), int(minutes))
	}
}
