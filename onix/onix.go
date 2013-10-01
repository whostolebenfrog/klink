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
    if args.SecondPos == "" {
        console.Fail("Must supply an application name as second positional argument")
    }

    createBody := Service{args.SecondPos}

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

func ServiceExists(serviceName string) bool {
    resp, _ := common.Head(onixUrl("/applications" + serviceName))
    return resp
}
