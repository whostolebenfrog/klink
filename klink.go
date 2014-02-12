package main

import (
	"fmt"
	optarg "github.com/jteeuwen/go-pkg-optarg"
	asgard "nokia.com/klink/asgard"
	common "nokia.com/klink/common"
	console "nokia.com/klink/console"
	ditto "nokia.com/klink/ditto"
	doctor "nokia.com/klink/doctor"
	exploud "nokia.com/klink/exploud"
	jenkins "nokia.com/klink/jenkins"
	onix "nokia.com/klink/onix"
	props "nokia.com/klink/props"
	tyr "nokia.com/klink/tyranitar"
	update "nokia.com/klink/update"
	"os"
	"strings"
)

var cmd = `[command] [application] [options]

[Commands]
    add-onix-prop       {application} -N property name -V json value
    allow-prod          {application} Allows the prod aws account access to the supplied application (or base/public)
    bake                {application} -v {version}
                        Bakes an AMI for {application} with version {version}
    create-app          {application} -E {email} -o {owner} -d {description}
                        Creates a new application. You probably want this when adding
                        a new service.
    deploy              {application} {environment} {ami}
                        Deploy the AMI {ami} for {application} to {environment}
    doctor              Test that everything is setup for klink to function
    info                {application} Return information about the application
    list-amis           {application} Lists the latest amis for the supplied application name
    list-apps           Lists the applications that exist (via exploud)
    list-apps-onix      Lists the applications that exist (in onix)
    list-apps-tyr       Lists the applications that exist (in tyranitar)
    register-app-onix   {application}
                        Creates a new application in onix only, useful for services
                        that won't be deployed using the cloud tooling
    register-app-tyr    {application}
                        Creates a new application in tyranitar only
    rollback            {application} {environment} rolls the application back to the last
                        successful deploy
    status              {application} Checks the status of the app
    undo                {application} {environment} Undo the steps of a broken deployment
    update              Update to the current version of klink.`

func printHelpAndExit() {
	console.Klink()
	console.Green()
	update.PrintVersion()
	console.Reset()
	fmt.Print("\n[New and updated] ")
	console.Red()
	fmt.Print("undo\n")
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
	optarg.Add("e", "environment", "Sets the environment", "poke, prod")
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
	case "rollback":
		exploud.Rollback(args)
	case "undo":
		exploud.Undo(args)
	case "bake":
		ditto.Bake(args)
	case "allow-prod":
		ditto.AllowProd(args)
	case "register-app-onix":
		onix.CreateApp(args)
	case "list-apps-onix":
		onix.ListApps()
	case "register-app-tyr":
		tyr.CreateApp(args)
	case "list-apps-tyr":
		tyr.ListApps()
	case "list-apps":
		exploud.ListApps()
	case "list-servers":
		asgard.ListServers(args)
	case "create-app":
		exploud.CreateApp(args)
	case "doctor":
		doctor.Doctor(args)
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
	case "ditto":
		ditto.Helpers(args)
	case "speak":
		console.Speak(args)
	case "build":
		jenkins.Build(args)
	default:
		printHelpAndExit()
	}
}

func main() {
	props.EnsureRCFile()
	update.EnsureUpdatedRecently(os.Args[0])
	handleAction(loadFlags())
}
