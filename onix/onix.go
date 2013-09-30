package onix

import(
    "fmt"
    common "nokia.com/klink/common"
    console "nokia.com/klink/console"
)

type Service struct {
    Name string `json:"name"`
}

func onixUrl(end string) string {
    return "http://onix.brislabs.com:8080/1.x" + end
}

func CreateService(args common.Command) {
    if args.Application == "" {
        console.Fail("Must supply an application name")
    }

    createBody := Service{args.Application}

    response, err := common.PostJson(onixUrl("/applications"), createBody)

    if err != nil {
        fmt.Println(err)
        console.BigFail("Unable to register new service with onix")
    }

    fmt.Println("Onix has created our service for us!")
    fmt.Println(response)
}

func ListServices() {
    response, err := common.GetString(onixUrl("/applications"))
    if err != nil {
        fmt.Println(err)
        console.Fail("Error listing services")
    }
    fmt.Println(response)
}
