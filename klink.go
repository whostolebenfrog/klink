package main

import (
	"fmt"
        "strings"
	optarg "github.com/jteeuwen/go-pkg-optarg"
	common "nokia.com/klink/common"
	console "nokia.com/klink/console"
	ditto "nokia.com/klink/ditto"
	exploud "nokia.com/klink/exploud"
	onix "nokia.com/klink/onix"
	tyr "nokia.com/klink/tyranitar"
	update "nokia.com/klink/update"
	"os"
)

var (
	cmd = "[Commands]\n" +
		"bake                {application} -v {version}\n" +
		"                    Bakes an AMI for {application} with version {version}.\n" +
		"create-service      {application} -E {email} -o {owner} -d {description}\n" +
		"                    Creates a new service via expload.\n" +
		"create-service-onix {application}\n" +
		"                    Creates a new service in onix only.\n" +
		"create-service-tyr  {application}\n" +
		"                    Creates a new service in tyranitar only.\n" +
		"deploy              {application} -a {ami}\n" +
		"                    Deploy the AMI {ami} for {application}.\n" +
		"doctor              Not yet implemented.\n" +
		"list-services       Lists the services that exist (in onix).\n" +
		"list-services-tyr   Lists the services that exist (in tyranitar).\n" +
		"update              Update to the current version of klink."
)

func printHelpAndExit() {
	console.Klink()
	update.PrintVersion()
	fmt.Println("\n")
	fmt.Println(strings.Replace(optarg.UsageString(), "[options]", "command [application] [options]\n\n" + cmd, 1))
	os.Exit(0)
}

// TODO: general - doc string on functions?
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
	case "create-service-onix":
		onix.CreateService(args)
	case "list-services":
		onix.ListServices()
	case "create-service-tyr":
		tyr.CreateService(args)
	case "list-services-tyr":
		tyr.ListServices()
	case "create-service":
		exploud.CreateService(args)
	case "doctor":
		fmt.Println("The Doctor is in the house")
	default:
		console.Fail(fmt.Sprintf("Unknown or no action: %s", args.Action))
	}
}

func main() {
	handleAction(loadFlags())
}
