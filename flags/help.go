package flags

import (
	"bytes"
	"fmt"
	optarg "github.com/jteeuwen/go-pkg-optarg"
	common "nokia.com/klink/common"
	console "nokia.com/klink/console"
	update "nokia.com/klink/update"
	"os"
	"strings"
	"text/tabwriter"
)

func PrintHelpAndExit() {
	// Top text
	console.Klink()
	console.Green()
	update.PrintVersion()
	console.Reset()
	fmt.Print("\n[New and updated] ")
	console.Red()
	fmt.Print("boxes, build, clone-tyr, clone-shuppet, delete-onix-prop\n")
	console.FReset()

	// [commands]
	w := new(tabwriter.Writer)
	output := new(bytes.Buffer)

	w.Init(output, 18, 8, 0, '\t', 0)
	for i := range common.Components {
		helpString := "\n    " + common.Components[i].String()
		fmt.Fprint(w, helpString)
	}
	w.Flush()
	outString := "[command] [app] [options]\n\n[commands]" + string(output.Bytes())

	// optstring and commands
	fmt.Println(strings.Replace(optarg.UsageString(), "[options]:", outString, 1))
	os.Exit(0)
}
