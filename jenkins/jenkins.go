package jenkins

import (
	"fmt"
	"net/http"
	common "nokia.com/klink/common"
	console "nokia.com/klink/console"
	onix "nokia.com/klink/onix"
	"strings"
	"time"
)

// Build a release job for the supplied application and poll the reponse
func Build(args common.Command) {
	app := args.SecondPos
	if app == "" {
		console.Fail("Yeah, you're gonna have to tell me what to build...")
	}

	path := BuildPath(app)
	fmt.Println("Calling release job at URL: " + path)

	location := PostBuild(path) + "api/json"
    job := GetJobFromQueue(location, 12)

    console.Green()
    fmt.Println("\nBuild started, polling jenkins for ouput...\n")
    console.Reset()

    PollBuild(job)
}

// Returns the release path for the supplied app
func BuildPath(app string) string {
	releasePath := onix.GetProperty(app, "releasePath")
	if !strings.HasSuffix(releasePath, "/") {
		releasePath += "/"
	}
	return releasePath + "build/api/json"
}

// Post a build and return the jobs location in the queue
func PostBuild(url string) string {
	resp, err := http.Post(url, "", nil)
	if err != nil {
		fmt.Println(fmt.Sprintf("Error posting build job to path: %s", url))
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
			fmt.Println("Unable to parse jenkins reponse, build may be in a queue")
			panic(err)
		}
	}

	return obj
}

// Poll a build and print the status
func PollBuild(path string) {
    status := GetJobStatus(path)
    lines := GetJobOutput(path)
    offset := 0
	timeout := time.Now().Add((20 * time.Minute))

    for (status == "in progress...") && time.Now().Before(timeout) {

        for i := offset; i < len(lines); i++ {
            fmt.Println(lines[i])
        }
        offset = len(lines)

        time.Sleep(1 * time.Second)
        status = GetJobStatus(path)
        lines = GetJobOutput(path)
    }
    console.Green()
    fmt.Println(status)
    console.Reset()
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

func GetJobOutput(path string) []string {
    path += "logText/progressiveText?start=0"
    return strings.Split(common.GetString(path), "\n")
}
