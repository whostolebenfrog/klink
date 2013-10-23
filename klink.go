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

var cmd = `[command] [application] [options]

[Commands]
    add-onix-prop       {application} -N property name -V json value
    bake                {application} -v {version}
                        Bakes an AMI for {application} with version {version}.
    create-app          {application} -E {email} -o {owner} -d {description}
                        Creates a new application via exploud.
    create-app-onix     {application}
                        Creates a new application in onix only.
    create-app-tyr      {application}
                        Creates a new application in tyranitar only.
    deploy              {application} {environment} {ami}
                        Deploy the AMI {ami} for {application} to {environment}.
    doctor              Not yet implemented.
    info                {application} Return information about the application
    list-apps           Lists the applications that exist (via exploud)
    list-apps-onix      Lists the applications that exist (in onix).
    list-apps-tyr       Lists the applications that exist (in tyranitar).
    list-amis           {application} Lists the latest amis for the supplied application name
    status              {application} Checks the status of the app
    update              Update to the current version of klink.`

func printHelpAndExit() {
	console.Klink()
	console.Green()
	update.PrintVersion()
	console.Reset()
	fmt.Print("\n[New and updated] ")
	console.Red()
	fmt.Print("status\n")
	console.FReset()
	fmt.Println(strings.Replace(optarg.UsageString(), "[options]:", cmd, 1))
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
	optarg.Add("e", "environment", "Sets the environment", "dev, prod")
	optarg.Add("m", "message", "Sets an informational message", "")
	optarg.Add("v", "version", "Sets the version", "")
	optarg.Add("d", "description", "Set the description for commands that require it", "")
	optarg.Add("E", "email", "Sets the email address for commands that require it", "")
	optarg.Add("o", "owner", "Sets the owner name for commands that require it", "")
	optarg.Add("s", "silent", "Sets silent mode, don't report to hubot", "")
	optarg.Add("D", "debug", "Sets debug mode. Gives more info on fails.", "")
	optarg.Add("N", "name", "Sets the property name", "")
	optarg.Add("V", "value", "Sets the property value", "")

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
		case "s":
			command.Silent = opt.Bool()
		case "D":
			command.Debug = opt.Bool()
		case "N":
			command.Name = opt.String()
		case "V":
			command.Value = opt.String()
		}
	}

	// positional arguments
	if len(os.Args) < 2 {
		printHelpAndExit()
	}
	command.Action = os.Args[1]
	// some commands need a second positional argument
	// let's do this better eh!?
	if len(os.Args) > 2 {
		command.SecondPos = os.Args[2]
	}
	if len(os.Args) > 3 {
		command.ThirdPos = os.Args[3]
	}
	if len(os.Args) > 4 {
		command.FourthPos = os.Args[4]
	}

	return command
}

func handleAction(args common.Command) {
	defer func() {
		if p := recover(); p != nil {
			if args.Debug == true {
				panic(p)
			}
			console.Red()
			fmt.Println(p)
			console.Reset()
			console.Fail("An error has occured. You may get more information using --debug true")
		}
	}()

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
	case "list-amis":
		ditto.FindAmis(args)
	case "find-amis":
		fmt.Println("Did you mean list-amis?")
	case "info":
		onix.Info(args)
	case "add-onix-prop":
		onix.AddProperty(args)
    case "status":
        onix.Status(args)
	default:
		printHelpAndExit()
	}
}

func main() {
	handleAction(loadFlags())
}
