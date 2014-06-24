package flags

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strings"
	"text/tabwriter"

	optarg "github.com/jteeuwen/go-pkg-optarg"
	common "nokia.com/klink/common"
	console "nokia.com/klink/console"
	update "nokia.com/klink/update"
)

func WrapString(desc string) []string {
	// heheheh - so much filth.
	re := regexp.MustCompile(".{1,80}(?:\\W|$)")
	return re.FindAllString(desc, -1)
}

func PrintHelpAndExit() {
	// Top text
	console.Klink()
	console.Green()
	update.PrintVersion()
	console.Reset()
	fmt.Print("\n[New and updated] ")
	console.Red()
	fmt.Print("bake retries, test, jobs\n")
	console.FReset()

	// [commands]
	w := new(tabwriter.Writer)
	output := new(bytes.Buffer)

	w.Init(output, 18, 8, 0, '\t', 0)
	for i := range common.Components {
		helpString := "\n    "
		desc := common.Components[i].String()

		// wrap long lines at white space
		if len(desc) > 80 {
			lines := WrapString(desc)
			for i := range lines {
				helpString += lines[i]
				if i < (len(lines) - 1) {
					helpString += "\n\t"
				}
			}
		} else {
			helpString += desc
		}

		fmt.Fprint(w, helpString)
	}
	w.Flush()
	outString := "[command] [app] [options]\n\n[commands]" + string(output.Bytes())

	// optstring and commands
	fmt.Println(strings.Replace(optarg.UsageString(), "[options]:", outString, 1))
	os.Exit(0)
}
