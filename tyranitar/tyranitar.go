package tyranitar

import (
    "fmt"
    common "nokia.com/klink/common"
    console "nokia.com/klink/console"
)

func tyrUrl(end string) {
    return "http://tyranitar.brislabs.com:8080/1.x" + end
}

func CreateServiceData(args common.Command) {
    if args.Application == "" {
        console.Fail("Must supply an application name")
    }
}

func ListServices() {
}
