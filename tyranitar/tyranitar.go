package tyranitar

import (
    "fmt"
    common "nokia.com/klink/common"
    console "nokia.com/klink/console"
)

type Service struct {
    Name string `json:"name"`
}

func tyrUrl(end string) string {
    return "http://tyranitar.brislabs.com:8080/1.x" + end
}

func CreateService(args common.Command) {
    if args.Application == "" {
        console.Fail("Must supply an application name")
    }

    createBody := Service{args.Application}

    response, err := common.PostJson(tyrUrl("/applications"), createBody)

    if err != nil {
        fmt.Println(err)
        console.BigFail("Unable to register new serivce with tyranitar")
    }

    fmt.Println("We are not registered with tyranitar!")
    fmt.Println(response)
}

func ListServices() {
    response, err := common.GetString(tyrUrl("/applications"))
    if err != nil {
        fmt.Println(err)
        console.Fail("Error listing services with Tyranitar")
    }
    fmt.Println(response)
}
