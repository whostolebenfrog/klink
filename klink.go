package main

import (
	"fmt"
	flags "nokia.com/klink/flags"
	common "nokia.com/klink/common"
	console "nokia.com/klink/console"
	ditto "nokia.com/klink/ditto"
	doctor "nokia.com/klink/doctor"
	exploud "nokia.com/klink/exploud"
	git "nokia.com/klink/git"
	jenkins "nokia.com/klink/jenkins"
	onix "nokia.com/klink/onix"
	props "nokia.com/klink/props"
	update "nokia.com/klink/update"
	"os"
)

func handleAction(args common.Command) {
    // global error handling
	defer func() {
		if p := recover(); p != nil {
			if args.Debug == true {
				console.Red()
                fmt.Println("\nDon't worry about the paths in trace, that's just go.\n")
				console.Reset()
				panic(p)
			}
			console.Red()
			fmt.Println(p)
			console.Reset()
			console.Fail("An error has occured. You may get more information using --debug true")
		}
	}()

    // special case here as requires os.Args not common.Command
	if args.Action == "update" {
		update.Update(os.Args[0])
		return
	}

    // everything else
	for i := range common.Components {
		component := common.Components[i]
		if args.Action == component.Command {
			component.Callback(args)
			return
		}
	}

    // failed to find the command, print help
	flags.PrintHelpAndExit()
}

func init() {
	// This whole thing makes me sad. Go demands that stuff like this is explicit
	// if we don't reference the namespace then even the .init() function won't be
	// called. We can't reference the namespace without using it so we basically
	// need to manually call the psuedo init methods, Init(), on each component
	// namesapce. Go doesn't allow, or encourage, this kind of aspecty metaprogramming
	console.Init()
	ditto.Init()
	doctor.Init()
	exploud.Init()
	git.Init()
	onix.Init()
	jenkins.Init()
}

func main() {
	props.EnsureRCFile()
	update.EnsureUpdatedRecently(os.Args[0])
	handleAction(flags.LoadFlags())
}
