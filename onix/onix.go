package onix

import(
    "fmt"
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

// Adds a property in onix to the supplied application
//PUT /1.x/applications/{app}/{key} (create/update a metadata item)
//    PUT http://onix.ent.nokia.com:8080/1.x/applications/empoleon/stats
//    {
//      "value" : [
//        {
//          "hp" : 84
//        },
//        {
//          "attack" : 86
//        },
//        {
//          "defense" : 88
//        }
//      ]
//    }
func AddProperty(args common.Command) {
}
