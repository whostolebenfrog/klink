package tyranitar

import (
    "fmt"
    common "nokia.com/klink/common"
    console "nokia.com/klink/console"
)

func CreateServiceData(args common.Command) {
    if args.Application == "" {
        console.Fail("Must supply an application name")
    }

}
