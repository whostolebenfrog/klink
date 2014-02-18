package git

import (
	"fmt"
    common "nokia.com/klink/common"
    console "nokia.com/klink/console"
    "os/exec"
)

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

func sncUrlTyranitar(app string, env string) string {
	return fmt.Sprintf("ssh://snc@source.nokia.com/tyranitar/git/%s-%s", app, env)
}

func sncUrlShuppet(app string) string {
	return fmt.Sprintf("ssh://snc@source.nokia.com/shuppet/git/%s", app)
}

func gitClone(path string) {
	out, err := exec.Command("git", "clone", path).Output()
    if err != nil {
        panic(err)
    }
    fmt.Println(string(out))
}

// Clone the tyranitar properties for the supplied app
func CloneTyranitar(args common.Command) {
	app := appName(args)
	env := envName(args)

	if env == "all" {
		gitClone(sncUrlTyranitar(app, "poke"))
		gitClone(sncUrlTyranitar(app, "prod"))
	} else {
		gitClone(sncUrlTyranitar(app, env))
	}
}

// Clone the shuppet properties for the supplied app
func CloneShuppet(args common.Command) {
	app := appName(args)
    gitClone(sncUrlShuppet(app))
}
