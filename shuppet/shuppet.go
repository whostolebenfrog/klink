package shuppet

import (
	"fmt"

	common "nokia.com/klink/common"
	console "nokia.com/klink/console"
)

func Init() {
	common.Register(
		common.Component{"show-infra", ShowInfra,
			"{app} {env} Shows infrastructure configuration for {app} in {env}", "APPS:ENVS"},
		common.Component{"apply-infra", ApplyInfra,
			"{app} {env} Apply infrastructure configuration for {app} in {env}", "APPS:ENVS"})
}

func shuppetUrl(end string) string {
	return "http://shuppet.brislabs.com" + end
}

// Returns all information stored in Shuppet about the supplied application and environment
func ShowInfra(args common.Command) {
	if args.SecondPos == "" {
		console.Fail("You must supply an app as the second positional argument")
	}
	app := args.SecondPos
	if args.ThirdPos == "" {
		console.Fail("You must supply an environment as the third positional argument")
	}
	env := args.ThirdPos

	infraUrl := fmt.Sprintf(shuppetUrl("/1.x/envs/%s/apps/%s"), env, app)
	console.MaybeJQS(common.GetString(infraUrl))
}

// Apply infrastructure configuration for supplied application and environment
func ApplyInfra(args common.Command) {
	if args.SecondPos == "" {
		console.Fail("You must supply an app as the second positional argument")
	}
	app := args.SecondPos
	if args.ThirdPos == "" {
		console.Fail("You must supply an environment as the third positional argument")
	}
	env := args.ThirdPos

	infraUrl := fmt.Sprintf(shuppetUrl("/1.x/envs/%s/apps/%s/apply"), env, app)
	console.MaybeJQS(common.GetString(infraUrl))
}
