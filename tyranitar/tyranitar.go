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

    response := common.PostJson(tyrUrl("/applications"), createBody)

    fmt.Println("Service registered with tyrnaitar")
    fmt.Println(response)
}

func ListApps() {
    fmt.Println(common.GetString(tyrUrl("/applications")))
}
