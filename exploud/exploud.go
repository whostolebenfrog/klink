package exploud

import (
	"encoding/json"
	"fmt"
	common "nokia.com/klink/common"
	console "nokia.com/klink/console"
//	"os"
	"strings"
//	"time"
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

    /*
	fmt.Println("Starting our rewriting tests")

    fmt.Println("Line 1 to overwrite")
    fmt.Println("Another line to overwrite")

    time.Sleep(500 * time.Millisecond)

	for i := 0; i < 10; i++ {
        time.Sleep(100 * time.Millisecond)
        fmt.Println("\033[31m\033[2A\rOverwriting first line", i)
        fmt.Println(fmt.Sprintf("\033[%dmOverwriting second line  ", 30 + i))
	}

    fmt.Print("\033[0m")

	os.Exit(0)
    */

	deployRequest := DeployRequest{args.Ami, "dev"}
	b, err := json.Marshal(deployRequest)
	if err != nil {
		console.Fail("Unable to create exploud deploy requset body")
	}
	fmt.Println("Calling exploud:", deployUrl, string(b))

    resp, err := common.PostJson(deployUrl, &deployRequest)
	if err != nil {
        fmt.Println(err, resp)
		console.Fail("Error calling exploud, exiting.")
	}

    fmt.Println(resp)
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
