package main

import (
	"flag"
	"fmt"
    "os"
	common "nokia.com/klink/common"
	console "nokia.com/klink/console"
	ditto "nokia.com/klink/ditto"
	exploud "nokia.com/klink/exploud"
	onix "nokia.com/klink/onix"
    tyr "nokia.com/klink/tyranitar"
    update "nokia.com/klink/update"
)

// TODO: positional args!
func loadFlags() common.Command {
	command := common.Command{}
	flag.StringVar(&command.Action, "action", "", "Action for klink to perform: [deploy, build, rollback]")
	flag.StringVar(&command.Application, "app", "", "Application to do stuff with")
	flag.StringVar(&command.Ami, "ami", "deploy", "Set the ami to deploy")
	flag.StringVar(&command.Environment, "env", "ent-dev", "Environment for this command [ent-dev, prod]")
	flag.StringVar(&command.Message, "m", "", "Message for this action. For example why you are deploying")
	flag.StringVar(&command.Version, "v", "", "Version to deploy, rollback etc [e.g. 0.153]")
	flag.Parse()

	fmt.Println("\t\t\tAction:      ", command.Action)
	fmt.Println("\t\t\tAmi:         ", command.Ami)
	fmt.Println("\t\t\tApplication: ", command.Application)
	fmt.Println("\t\t\tEnvironment: ", command.Environment)
	fmt.Println("\t\t\tMessage:     ", command.Message)
	fmt.Println("\t\t\tVersion:     ", command.Version, "\n")
	return command
}

// TODO: figure out some names here
// TODO: --json output mode
// TODO: DOCTOR! 
func handleAction(args common.Command) {
	switch args.Action {
    case "version":
        update.Version()
    case "update":
        update.Update(os.Args[0])
	case "deploy":
		exploud.Exploud(args)
	case "bake":
		ditto.Bake(args)
	case "create-service-onix":
		onix.CreateService(args)
    case "list-services-onix":
        onix.ListServices()
    case "create-service-tyr":
        tyr.CreateService(args)
    case "list-services-tyr":
        tyr.ListServices()
    case "create-service":
        onix.CreateService(args)
        tyr.CreateService(args)
	default:
		console.Fail(fmt.Sprintf("Unknown or no action: %s", args.Action))
	}
}

func main() {
	console.Klink()
	command := loadFlags()
	handleAction(command)
}
