package onix

import (
	"encoding/json"
	"fmt"
	jsonq "github.com/jmoiron/jsonq"
	common "nokia.com/klink/common"
	console "nokia.com/klink/console"
)

func Init() {
	common.Register(
		common.Component{"register-app-onix", CreateApp,
			"{app} Creates a new application in onix only"},
		common.Component{"list-apps-onix", ListApps,
			"Lists the applications that exist in onix"},
		common.Component{"info", Info,
			"{app} Return information about the application"},
		common.Component{"add-onix-prop", AddProperty,
			"{app} -N property name -V json value"},
		common.Component{"get-onix-prop", GetPropertyFromArgs,
			"{app} {property-name} get the property for the application"},
		common.Component{"status", Status,
			"{app} Checks the status of the app"},
		common.Component{"delete-onix-prop", DeleteProperty,
			"{app} {property-name} Delete the property for the application"})
}

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
func ListApps(args common.Command) {
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
	if args.SecondPos == "" {
		console.Fail("Must supply service name as a second positional argument")
	}
	if args.Name == "" {
		console.Fail("Must supply property name using -N")
	}
	value := args.Value
	if value == "" {
		console.Fail("Must supply value using -V in json format. Remember to quote!")
	}

	valueString, err := ToJsonValue(value)
	if err != nil {
		valueString, err = ToJsonValue("\"" + value + "\"")
		if err != nil {
			fmt.Println("That doesn't look like json")
			panic(err)
		}
	}

	fmt.Println(common.PutString(onixUrl("/applications/"+args.SecondPos+"/"+args.Name),
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
			console.Fail(fmt.Sprintf("klink add-onix-prop %s -N %s -V '[\"my\", \"array\"]'\n",
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
		console.Fail("Must supply service name as a second positional argument")
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

func GetPropertyFromArgs(args common.Command) {
	app := args.SecondPos
	name := args.ThirdPos
	if app == "" {
		console.Fail("Don't forget to bring a towel ^H^H^H^H pass a application name")
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
