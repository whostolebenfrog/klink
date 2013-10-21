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

func CreateApp(args common.Command) {
    if args.SecondPos == "" {
        console.Fail("Must supply an application name as second positional argument")
    }

    createBody := App{args.SecondPos}

    response := common.PostJson(onixUrl("/applications"), createBody)

    fmt.Println("Onix has created our application for us!")
    fmt.Println(response)
}

func ListApps() {
    fmt.Println(common.GetString(onixUrl("/applications")))
}

func AppExists(appName string) bool {
    return common.Head(onixUrl("/applications/" + appName))
}
