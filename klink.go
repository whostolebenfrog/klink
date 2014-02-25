package main

import (
    "bytes"
	"fmt"
	optarg "github.com/jteeuwen/go-pkg-optarg"
	asgard "nokia.com/klink/asgard"
	common "nokia.com/klink/common"
	console "nokia.com/klink/console"
    ditto "nokia.com/klink/ditto"
	doctor "nokia.com/klink/doctor"
	exploud "nokia.com/klink/exploud"
    git "nokia.com/klink/git"
	jenkins "nokia.com/klink/jenkins"
	onix "nokia.com/klink/onix"
	props "nokia.com/klink/props"
	tyr "nokia.com/klink/tyranitar"
	update "nokia.com/klink/update"
	"os"
	"strings"
    "text/tabwriter"
)

var cmd = `[command] [application] [options]

[Commands]
    add-onix-prop       {application} -N property name -V json value
    build               {application} builds the jenkins release job for an application
    clone-tyr           {application} {env} clone the tyranitar properties for an app. Pass {env}
                        to optionally only clone that env. Defaults to all
    clone-shuppet       {application} {env} clone the shuppet properties for an app. Pass {env}
                        to optionally only clone that env. Defaults to all
    doctor              Test that everything is setup for klink to function
    get-onix-prop       {application} {property-name} get the property for the application
    info                {application} Return information about the application
    list-apps-onix      Lists the applications that exist (in onix)
    list-apps-tyr       Lists the applications that exist (in tyranitar)
    register-app-onix   {application}
                        Creates a new application in onix only, useful for services
                        that won't be deployed using the cloud tooling
    register-app-tyr    {application}
                        Creates a new application in tyranitar only
    status              {application} Checks the status of the app
    update              Update to the current version of klink.`

func printHelpAndExit() {
    // Top text
	console.Klink()
	console.Green()
	update.PrintVersion()
	console.Reset()
	fmt.Print("\n[New and updated] ")
	console.Red()
	fmt.Print("boxes, build, clone-tyr, clone-shuppet\n")
	console.FReset()

    // [commands]
    w := new(tabwriter.Writer)
    output := new(bytes.Buffer)

	w.Init(output, 18, 8, 0, '\t', 0)
    for i := range(common.Components) {
        helpString := "\n    " + common.Components[i].String()
        fmt.Fprint(w, helpString)
    }
	w.Flush()
    outString := "[command] [app] [options]\n\n[commands]" + string(output.Bytes())

    // optstring and commands
	fmt.Println(strings.Replace(optarg.UsageString(), "[options]:", outString, 1))
	os.Exit(0)
}

func loadFlags() common.Command {
	command := common.Command{}

	// flags
	optarg.Header("General Options")
	optarg.Add("h", "help", "Displays this help message", false)
	optarg.Header("Deployment based flags")
	optarg.Add("a", "ami", "Sets the ami for commands that require it", "")
	optarg.Add("d", "description", "Set the description for commands that require it", "")
	optarg.Add("D", "debug", "Sets debug mode. Gives more info on fails.", "")
	optarg.Add("e", "environment", "Sets the environment", "poke, prod")
	optarg.Add("E", "email", "Sets the email address for commands that require it", "")
	optarg.Add("f", "format", "Sets the format property value", "")
	optarg.Add("m", "message", "Sets an informational message", "")
	optarg.Add("N", "name", "Sets the property name", "")
	optarg.Add("o", "owner", "Sets the owner name for commands that require it", "")
	optarg.Add("s", "silent", "Sets silent mode, don't report to hubot", "")
	optarg.Add("S", "status", "Sets the status property value", "")
	optarg.Add("v", "version", "Sets the version", "")
	optarg.Add("V", "value", "Sets the property value", "")

	for opt := range optarg.Parse() {
		switch opt.ShortName {
		case "a":
			command.Ami = opt.String()
		case "d":
			command.Description = opt.String()
		case "D":
			command.Debug = opt.Bool()
		case "e":
			command.Environment = opt.String()
		case "E":
			command.Email = opt.String()
        case "f":
            command.Format = opt.String()
		case "h":
			printHelpAndExit()
		case "m":
			command.Message = opt.String()
		case "N":
			command.Name = opt.String()
		case "o":
			command.Owner = opt.String()
		case "s":
			command.Silent = opt.Bool()
        case "S":
            command.Status = opt.String()
		case "v":
			command.Version = opt.String()
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
                console.Red()
                fmt.Println("\nThis may print out a trace with bgriffit's home directory - don't worry that's just where it was built.")
                console.Reset()
				panic(p)
			}
			console.Red()
			fmt.Println(p)
			console.Reset()
			console.Fail("An error has occured. You may get more information using --debug true")
		}
	}()

    for i := range common.Components {
        component := common.Components[i]
        if args.Action == component.Command {
            component.Callback(args)
            return
        }
    }

	switch args.Action {
	case "update":
		update.Update(os.Args[0])
	case "register-app-onix":
		onix.CreateApp(args)
	case "list-apps-onix":
		onix.ListApps()
	case "register-app-tyr":
		tyr.CreateApp(args)
	case "list-apps-tyr":
		tyr.ListApps()
	case "list-servers":
		asgard.ListServers(args)
	case "doctor":
		doctor.Doctor(args)
	case "info":
		onix.Info(args)
	case "add-onix-prop":
		onix.AddProperty(args)
    case "get-onix-prop":
        onix.GetPropertyFromArgs(args)
	case "status":
		onix.Status(args)
	case "speak":
		console.Speak(args)
	case "build":
		jenkins.Build(args)
    case "clone-tyr":
        git.CloneTyranitar(args)
    case "clone-shuppet":
        git.CloneShuppet(args)
	default:
		printHelpAndExit()
	}
}

func init() {
    // This whole thing makes me sad. Go demands that stuff like this is explicit
    // if we don't reference the namespace then even the .init() function won't be
    // called. We can't reference the namespace without using it so we basically
    // need to manually call the psuedo init methods, Init(), on each component
    // namesapce. Go doesn't allow, or encourage, this kind of aspecty metaprogramming
    ditto.Init()
    exploud.Init()
}

func main() {
	props.EnsureRCFile()
	update.EnsureUpdatedRecently(os.Args[0])
	handleAction(loadFlags())
}

