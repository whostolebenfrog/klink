package onix

import (
	"fmt"
	jsonq "github.com/jmoiron/jsonq"
	common "nokia.com/klink/common"
	console "nokia.com/klink/console"
)

type App struct {
	Name string `json:"name"`
}

func onixUrl(end string) string {
	return "http://onix.brislabs.com:8080/1.x" + end
}

// Create a new application in onix
func CreateApp(args common.Command) {
	if args.SecondPos == "" {
		console.Fail("Must supply an application name as second positional argument")
	}

	createBody := App{args.SecondPos}

	response := common.PostJson(onixUrl("/applications"), createBody)

	fmt.Println("Onix has created our application for us!")
	fmt.Println(response)
}

// List all apps that onix knows about
func ListApps() {
	fmt.Println(common.GetString(onixUrl("/applications")))
}

// Returns true if the app exists
func AppExists(appName string) bool {
	return common.Head(onixUrl("/applications/" + appName))
}

// Returns all information stored in onix about the supplied application
func Info(args common.Command) {
	fmt.Println(common.GetString(onixUrl("/applications/" + args.SecondPos)))
}

func AddProperty(args common.Command) {
	if args.SecondPos == "" {
		console.Fail("Must supply service name as a second positional argument")
	}
	if args.Name == "" {
		console.Fail("Must supply property name using -N")
	}
	if args.Value == "" {
		console.Fail("Must supply value using -V in json format. Remember to quote!")
	}
	fmt.Println(common.PutString(onixUrl("/applications/"+args.SecondPos+"/"+args.Name),
		args.Value))
}

func EnsureProp(jq *jsonq.JsonQuery, app string, name string) string {
	obj, err := jq.String("metadata", name)
	if err != nil {
		fmt.Println(fmt.Sprintf("Application %s doesn't have a %s defined, add one with:\n",
			app, name))
		console.Fail(fmt.Sprintf("klink add-onix-prop %s -N %s -V '{\"value\" : \"value\"}'\n",
			app, name))
	}
	return obj
}

func Status(args common.Command) {
	app := args.SecondPos
	if app == "" {
		console.Fail("Must supply service name as a second positional argument")
	}

	jq := common.GetAsJsonq(onixUrl("/applications/" + app))

	statusUrl := EnsureProp(jq, app, "servicePathPoke") + EnsureProp(jq, app, "statusPath")
	fmt.Println(fmt.Sprintf("Checking status at: %s", statusUrl))

	console.Green()
	fmt.Println(common.GetString(statusUrl))
	console.Reset()
}

func GetProperty(app string, name string) string {
	jq := common.GetAsJsonq(onixUrl("/applications/" + app))
	return EnsureProp(jq, app, name)
}
