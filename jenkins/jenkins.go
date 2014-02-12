package jenkins

import (
	common "nokia.com/klink/common"
    console "nokia.com/klink/console"
)

func Build(args common.Command) {
    if args.SecondPos == "" {
        console.Fail("Yeah, you're gonna have to tell me what to build...")
    }

}
