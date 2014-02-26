package flags

import (
	optarg "github.com/jteeuwen/go-pkg-optarg"
	common "nokia.com/klink/common"
	"os"
)

func LoadFlags() common.Command {
	command := common.Command{}

	// flags
	optarg.Header("General Options")
	optarg.Add("h", "help", "Displays this help message", false)
	optarg.Header("Deployment based flags")
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
			PrintHelpAndExit()
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
		PrintHelpAndExit()
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
