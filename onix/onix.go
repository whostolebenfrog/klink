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

    response, err := common.PostJson(onixUrl("/applications"), createBody)

    if err != nil {
        fmt.Println(err)
        console.BigFail("Unable to register new application with onix")
    }

    fmt.Println("Onix has created our application for us!")
    fmt.Println(response)
}

func ListApps() {
    response, err := common.GetString(onixUrl("/applications"))
    if err != nil {
        fmt.Println(err)
        console.Fail("Error listing applications")
    }
    fmt.Println(response)
}

func AppExists(appName string) bool {
    resp, _ := common.Head(onixUrl("/applications/" + appName))
    return resp
}
