package onix

import(
    "bytes"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    args "nokia.com/klink/args"
    console "nokia.com/klink/console"
)

type Service struct {
    Name string `json:"name"`
}

// TODO: abstract out all this read body stuff and http handling stuff!
func CreateService(args args.Command) {
    createServiceUrl := "http://onix.brislabs.com:8080/1.x/applications"

    createBody := Service{args.Application}
    b, err := json.Marshal(createBody)
    if err != nil {
        console.BigFail("Unable to create onix application request body")
    }

    fmt.Println("Calling onix to create service:", args.Application)

    resp, err := http.Post(createServiceUrl, "application/json", bytes.NewReader(b))
    if err != nil {
        console.BigFail("Error attempting to talk to onix. Screw you onix.")
    }
    defer resp.Body.Close()

    if resp.StatusCode == 201 {
        body, err := ioutil.ReadAll(resp.Body)
        if err != nil {
            console.BigFail("Failed to read onix response body... WHY?!")
        }
        fmt.Println("Onix has created our service for us!")
        fmt.Println(string(body))
    } else {
        fmt.Println("Got non-201 response from onix")
        body, err := ioutil.ReadAll(resp.Body)
        if err != nil {
            console.BigFail("Failed to read onix response body... WHY?!")
        }
        fmt.Println((string(body)))
        console.BigFail("Big fail trying to create a service in onix :-(")
    }
}

func ListServices() {
    
}
