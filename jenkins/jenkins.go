package jenkins

import (
	"fmt"
	common "nokia.com/klink/common"
	console "nokia.com/klink/console"
	onix "nokia.com/klink/onix"
)

func Build(args common.Command) {
	app := args.SecondPos
	if app == "" {
		console.Fail("Yeah, you're gonna have to tell me what to build...")
	}
	fmt.Println(onix.GetProperty(app, "releasePath"))
}
