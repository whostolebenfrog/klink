package git

import (
	"fmt"
	"os/exec"

	common "mixrad.io/klink/common"
	console "mixrad.io/klink/console"
	onix "mixrad.io/klink/onix"
)

func Init() {
	common.Register(
		common.Component{"clone-tyr", CloneTyranitar,
			"{app} {env - optional} clone the tyranitar properties for an app into pwd", "APPS"},
		common.Component{"clone-shuppet", CloneShuppet,
			"{app} {env - optional} clone the shuppet properties for an app into pwd", "APPS|ENVS"},
		common.Component{"clone", CloneService,
			"{app} clone the application into pwd", "APPS"},
		common.Component{"gist", Gist,
			"{file-name} [{description}] send stdin to a github gist, use extension to set type", ""})
}

func appName(args common.Command) string {
	if args.SecondPos == "" {
		console.Fail("Application must be provided as the second positional argument")
	}
	return args.SecondPos
}

func envName(args common.Command) string {
	if args.ThirdPos == "" {
		return "all"
	}
	return args.ThirdPos
}

func gitUrlTyranitar(app string, env string) string {
	return fmt.Sprintf("git@github.brislabs.com:tyranitar/%s-%s.git", app, env)
}

func gitUrlShuppet(app string) string {
	return fmt.Sprintf("git@github.brislabs.com:shuppet/%s.git", app)
}

func gitClone(path string) {
	out, err := exec.Command("git", "clone", path).Output()
	if err != nil {
		fmt.Println(fmt.Sprintf("Error cloning repo, %s, does it already exist? %s", path, err))
	}
	fmt.Println(string(out))
}

// Clone the tyranitar properties for the supplied app
func CloneTyranitar(args common.Command) {
	app := appName(args)
	env := envName(args)

	if env == "all" {
		gitClone(gitUrlTyranitar(app, "poke"))
		gitClone(gitUrlTyranitar(app, "prod"))
	} else {
		gitClone(gitUrlTyranitar(app, env))
	}
}

// Clone the shuppet properties for the supplied app
func CloneShuppet(args common.Command) {
	app := appName(args)
	gitClone(gitUrlShuppet(app))
}

func CloneService(args common.Command) {
	app := args.SecondPos

	path := onix.GetProperty(app, "srcRepo")
	gitClone(path)
}
