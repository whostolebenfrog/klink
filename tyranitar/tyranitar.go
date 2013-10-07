package tyranitar

import (
    "fmt"
    common "nokia.com/klink/common"
    console "nokia.com/klink/console"
)

type App struct {
    Name string `json:"name"`
}

func tyrUrl(end string) string {
    return "http://tyranitar.brislabs.com:8080/1.x" + end
}

func CreateApp(args common.Command) {
    if args.SecondPos == "" {
        console.Fail("Must supply an application name as the second positional argument")
    }

    createBody := App{args.SecondPos}

    response, err := common.PostJson(tyrUrl("/applications"), createBody)

    if err != nil {
        fmt.Println(err)
        console.BigFail("Unable to register new application with tyranitar")
    }

    fmt.Println("We are not registered with tyranitar!")
    fmt.Println(response)
}

func ListApps() {
    response, err := common.GetString(tyrUrl("/applications"))
    if err != nil {
        fmt.Println(err)
        console.Fail("Error listing applications with Tyranitar")
    }
    fmt.Println(response)
}
