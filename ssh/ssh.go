package ssh

import (
	"fmt"
	jsonq "github.com/jmoiron/jsonq"
	common "nokia.com/klink/common"
	console "nokia.com/klink/console"
	exploud "nokia.com/klink/exploud"
	onix "nokia.com/klink/onix"
	props "nokia.com/klink/props"
	"os"
	"os/exec"
)

func Init() {
	common.Register(
		common.Component{"ssh", SSH, "{app} {env} [{numel-id}] SSH onto a server [-v true]. Uses SSHUsername from klink.rc if set"})
}

func SSH(args common.Command) {
	app := args.SecondPos
	env := args.ThirdPos
	id := args.FourthPos
	verbose := args.Version

	if app == "" {
		console.Fail("You must supply an app as the second positional argument")
	}
	if !onix.KnownEnvironment(env) {
		fmt.Println(fmt.Sprintf("env not supplied or incorrect, setting to poke. %s",
			props.GetEnvironments()))
		env = "poke"
	}
	if id == "" {
		id = "any"
	}

	boxesArray := make([]interface{}, 20)
	exploud.JsonBoxes(app, env, boxesArray)
	var ip string
	var numelId string

	for _, jsonBox := range boxesArray {
		if jsonBox == nil {
			break
		}
		jqBox := jsonq.NewQuery(jsonBox)
		numelId, _ = jqBox.String("numel-id")
		if id == numelId || id == "any" {
			ip, _ = jqBox.String("private-ip")
			break
		}
	}

	if ip == "" {
		fmt.Println("Unable to find a matching server, found (ignore the nils):")
		console.Fail(fmt.Sprintf("%s", boxesArray))
	} else {
		fmt.Println(fmt.Sprintf("About to connect to %s with ip %s", numelId, ip))
		writeSSHScript(ip, verbose != "")
	}
}

func writeSSHScript(ip string, verbose bool) {
	if common.IsWindows() {
		console.Fail("Can't ssh on windows. Well, klink can't anyway :-/ Talk to Ben if you really want this")
	}

	var sshargs []string
	if verbose {
		sshargs = append(sshargs, "-v")
	}
	if props.Get("SSHUsername") != "" {
		sshargs = append(sshargs, "-l", props.Get("SSHUsername"))
	}
	sshargs = append(sshargs, ip)

	cmd := exec.Command("ssh", sshargs...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		panic(err)
	}
}
