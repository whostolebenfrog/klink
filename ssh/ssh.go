package ssh

import (
	"fmt"
	"os"
	"os/exec"
    "strconv"
    "strings"
	"text/tabwriter"

	jsonq "github.com/jmoiron/jsonq"
	common "nokia.com/klink/common"
	console "nokia.com/klink/console"
	exploud "nokia.com/klink/exploud"
	onix "nokia.com/klink/onix"
	props "nokia.com/klink/props"
)

func Init() {
	common.Register(
		common.Component{"ssh", SSH, "{app} {env} [{numel-id}] SSH onto a server [-v true]. Uses SSHUsername from klink.rc if set", "APPS:ENVS"})
}

func chooseSSH(boxes []interface{}) string {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 8, 8, 0, '\t', 0)

	for i, jsonBox := range boxes {
		if jsonBox == nil {
			break
		}

		jqBox := jsonq.NewQuery(jsonBox)
        name, _ := jqBox.String("name")
        instanceId, _ := jqBox.String("instance-id")
        imageId, _ := jqBox.String("image-id")
        privateIp, _ := jqBox.String("private-ip")
        numelId, _ := jqBox.String("numel-id")
        opts := []string{name, instanceId, imageId, numelId, privateIp, strconv.Itoa(i)}
        outStr := strings.Join(opts, "\t") + "\n"
		fmt.Fprint(w, outStr)
	}
	w.Flush()

    choice, err := strconv.ParseInt(console.GetPrompt("Pick an instance to log in to:"), 10, 8)
    if boxes[choice] == nil || err != nil {
        console.Red()
        console.Fail("You failed at this simple task. Have you considered a career in management?")
        console.Reset()
    }
    ip, _ := jsonq.NewQuery(boxes[choice]).String("private-ip")

    return ip
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

	boxesArray := make([]interface{}, 20)
	exploud.JsonBoxes(app, env, boxesArray)

    if id == "" {
        ip := chooseSSH(boxesArray)
		fmt.Println(fmt.Sprintf("About to connect to %s", ip))
        doSomeSSH(ip, verbose != "")
    } else {
        var ip string
        var numelId string

        for _, jsonBox := range boxesArray {
            if jsonBox == nil {
                break
            }
            jqBox := jsonq.NewQuery(jsonBox)
            numelId, _ = jqBox.String("numel-id")
            if id == numelId || id == "" {
                ip, _ = jqBox.String("private-ip")
                break
            }
        }

        if ip == "" {
            fmt.Println("Unable to find a matching server, found (ignore the nils):")
            console.Fail(fmt.Sprintf("%s", boxesArray))
        } else {
            fmt.Println(fmt.Sprintf("About to connect to %s with ip %s", numelId, ip))
            doSomeSSH(ip, verbose != "")
        }
    }

}

func doSomeSSH(ip string, verbose bool) {
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
