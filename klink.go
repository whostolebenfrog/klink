package main

import (
	"fmt"
	optarg "github.com/jteeuwen/go-pkg-optarg"
	common "nokia.com/klink/common"
	console "nokia.com/klink/console"
	ditto "nokia.com/klink/ditto"
	exploud "nokia.com/klink/exploud"
	onix "nokia.com/klink/onix"
	tyr "nokia.com/klink/tyranitar"
	update "nokia.com/klink/update"
	"os"
	"strings"
)

var (
	cmd = "[Commands]\n" +
		"bake                {application} -v {version}\n" +
		"                    Bakes an AMI for {application} with version {version}.\n" +
		"create-app          {application} -E {email} -o {owner} -d {description}\n" +
		"                    Creates a new application via exploud.\n" +
		"create-app-onix     {application}\n" +
		"                    Creates a new application in onix only.\n" +
		"create-app-tyr      {application}\n" +
		"                    Creates a new application in tyranitar only.\n" +
		"deploy              {application} -a {ami}\n" +
		"                    Deploy the AMI {ami} for {application}.\n" +
		"doctor              Not yet implemented.\n" +
		"list-apps           Lists the applications that exist (via exploud)\n" +
		"list-apps-onix      Lists the applications that exist (in onix).\n" +
		"list-apps-tyr       Lists the applications that exist (in tyranitar).\n" +
		"update              Update to the current version of klink."
)

func printHelpAndExit() {
	console.Klink()
	update.PrintVersion()
	fmt.Println("\n")
	fmt.Println(strings.Replace(optarg.UsageString(), "[options]:", "command [application] [options]\n\n"+cmd, 1))
	os.Exit(0)
}

// TODO: general - json output mode? jq mode?
func loadFlags() common.Command {
	command := common.Command{}

	// flags
	optarg.Header("General Options")
	optarg.Add("h", "help", "Displays this help message", false)
	optarg.Header("Deployment based flags")
	optarg.Add("a", "ami", "Sets the ami for commands that require it", "")
	optarg.Add("e", "environment", "Sets the environment", "ent-dev")
	optarg.Add("m", "message", "Sets an informational message", "")
	optarg.Add("v", "version", "Sets the version", "")
	optarg.Add("d", "description", "Set the description for commands that require it", "")
	optarg.Add("E", "email", "Sets the email address for commands that require it", "")
	optarg.Add("o", "owner", "Sets the owner name for commands that require it", "")

	for opt := range optarg.Parse() {
		switch opt.ShortName {
		case "h":
			printHelpAndExit()
		case "a":
			command.Ami = opt.String()
		case "e":
			command.Environment = opt.String()
		case "m":
			command.Message = opt.String()
		case "v":
			command.Version = opt.String()
		case "d":
			command.Description = opt.String()
		case "E":
			command.Email = opt.String()
		case "o":
			command.Owner = opt.String()
		}
	}

	// positional arguments
	if len(os.Args) < 2 {
		printHelpAndExit()
	}
	command.Action = os.Args[1]
	// some commands need a second positional argument
	if len(os.Args) > 2 {
		command.SecondPos = os.Args[2]
	}

	return command
}

func handleAction(args common.Command) {
	switch args.Action {
	case "update":
		update.Update(os.Args[0])
	case "deploy":
		exploud.Exploud(args)
	case "bake":
		ditto.Bake(args)
	case "create-app-onix":
		onix.CreateApp(args)
	case "list-apps-onix":
		onix.ListApps()
	case "create-app-tyr":
		tyr.CreateApp(args)
	case "list-apps-tyr":
		tyr.ListApps()
	case "list-apps":
		exploud.ListApps()
	case "create-app":
		exploud.CreateApp(args)
	case "doctor":
		fmt.Println("The Doctor is in the house")
	default:
		console.Fail(fmt.Sprintf("Unknown or no action: %s", args.Action))
	}
}

func main() {
	handleAction(loadFlags())
}
