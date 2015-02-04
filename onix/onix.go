package onix

import (
	"encoding/json"
	"fmt"
	"strings"

	jsonq "github.com/jmoiron/jsonq"
	common "nokia.com/klink/common"
	conf "nokia.com/klink/conf"
	console "nokia.com/klink/console"
	props "nokia.com/klink/props"
)

func Init() {
	common.Register(
		common.Component{"register-app-onix", CreateApp,
			"{app} Creates a new application in onix only", "APPS"},
		common.Component{"info", Info,
			"{app} Return information about the application", "APPS"},
		common.Component{"add-onix-prop", AddProperty,
			"{app} {name} {value} Adds an onix property (json)", "APPS|PROPNAMES"},
		common.Component{"get-onix-prop", GetPropertyFromArgs,
			"{app} {property-name} get the property for the application", "APPS"},
		common.Component{"status", Status,
			"{app} Checks the status of the app", "APPS"},
		common.Component{"apps", ListApps,
			"Lists the applications that exist (via exploud)", "APPS"},
		common.Component{"delete-onix-prop", DeleteProperty,
			"{app} {property-name} Delete the property for the application", "APPS"})
}

type App struct {
	Name string `json:"name"`
}

func onixUrl(end string) string {
	return conf.ListerUrl + end
}

// Return the list of apps that are known about by onix
func GetApps() []string {
	apps, err := common.GetAsJsonq(onixUrl("/applications")).ArrayOfStrings("applications")
	if err != nil {
		panic(err)
	}
	return apps
}

func GetCommonPropertyNames() []string {
	return []string{"bakeType", "baker", "customBakeCommands", "jobsPath", "releasePath", "srcRepo", "servicePathPoke", "statusPath", "testPath"}
}

// List the apps known about by onix
func ListApps(args common.Command) {
	fmt.Println(strings.Join(GetApps(), "\n"))
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

// Returns true if the app exists
func AppExists(appName string) bool {
	return common.Head(onixUrl("/applications/" + appName))
}

// Returns all information stored in onix about the supplied application
func Info(args common.Command) {
	console.MaybeJQS(common.GetString(onixUrl("/applications/" + args.SecondPos)))
}

func ToJsonValue(in string) (string, error) {
	in = "{\"value\" : " + in + "}"
	var generic interface{}
	err := json.Unmarshal([]byte(in), &generic)
	if err != nil {
		return "", err
	}

	x, err := json.Marshal(generic)
	return string(x), err
}

func AddProperty(args common.Command) {
	app := args.SecondPos
	if app == "" {
		console.Fail("Must supply application name as a second positional argument")
	}
	name := args.ThirdPos
	if name == "" {
		console.Fail("Must supply the property name as the third positional argument")
	}
	value := args.FourthPos
	if value == "" {
		console.Fail("Must supply the property value as the fourth positional argument")
	}

	valueString, err := ToJsonValue(value)
	if err != nil {
		valueString, err = ToJsonValue("\"" + value + "\"")
		if err != nil {
			fmt.Println("That doesn't look like json")
			panic(err)
		}
	}

	fmt.Println(common.PutString(onixUrl("/applications/"+app+"/"+name),
		valueString))
}

func EnsureProp(jq *jsonq.JsonQuery, app string, name string) string {
	str, err := jq.String("metadata", name)
	if err != nil {
		obj, err := jq.Interface("metadata", name)
		if err != nil {
			fmt.Printf(
				"Application %s doesn't have a %s defined, add one with:\n",
				app,
				name,
			)
			console.Fail(fmt.Sprintf("klink add-onix-prop %s %s 'value'\n",
				app, name))
		}
		// this is the only way to get a string from an arbitary type in go...
		return fmt.Sprintf("%s", obj)
	}
	return str
}

func Status(args common.Command) {
	app := args.SecondPos
	if app == "" {
		console.Fail("Must supply application name as a second positional argument")
	}

	jq := common.GetAsJsonq(onixUrl("/applications/" + app))

	statusUrl := EnsureProp(jq, app, "servicePathPoke") + EnsureProp(jq, app, "statusPath")
	fmt.Printf("Checking status at: %s\n", statusUrl)

	console.Green()
	fmt.Println(common.GetString(statusUrl))
	console.Reset()
}

func GetProperty(app string, name string) string {
	jq := common.GetAsJsonq(onixUrl("/applications/" + app))
	return EnsureProp(jq, app, name)
}

func GetOptionalProperty(app string, name string) string {
	jq := common.GetAsJsonq(onixUrl("/applications/" + app))
	str, err := jq.String("metadata", name)
	if err != nil {
		obj, err := jq.Interface("metadata", name)
		if err != nil {
			return ""
		}
		// this is the only way to get a string from an arbitary type in go...
		return fmt.Sprintf("%s", obj)
	}
	return str
}

func GetPropertyFromArgs(args common.Command) {
	app := args.SecondPos
	name := args.ThirdPos
	if app == "" {
		console.Fail("Don't forget to bring a towel^H^H^H^H^H^H pass a application name")
	}
	if name == "" {
		console.Fail("You forgot to pass the property name")
	}
	fmt.Println(GetProperty(app, name))
}

func DeleteProperty(args common.Command) {
	app := args.SecondPos
	name := args.ThirdPos

	if app == "" {
		console.Fail("You forgot to pass the app name")
	}
	if name == "" {
		console.Fail("You forgot to pass the property name")
	}

	common.Delete(onixUrl("/applications/" + app + "/" + name))

	console.Green()
	fmt.Println("Success!")
	console.Reset()
}

// Get the list of environments from onix
func EnvironmentsFromOnix() []string {
	envs, err := common.GetAsJsonq(onixUrl("/environments")).ArrayOfStrings("environments")
	if err != nil {
		panic("Unable to parse response getting environments :-(")
	}
	return envs
}

// Returns a list of available environments, accepts an environment
// if that environment isn't known then go and ge the list from
// onix
func GetEnvironments(env string) []string {
	environments := props.GetEnvironments()
	if !common.Contains(environments, env) {
		environments = EnvironmentsFromOnix()
		props.SetEnvironments(environments)
	}
	return environments
}

// Returns true if the environment is known by onix
func KnownEnvironment(env string) bool {
	return common.Contains(GetEnvironments(env), env)
}
