package console

import (
	"fmt"
    "os"
	"os/exec"
)

// Sprint the supplied string or if jq is available use that to parse it
func MaybeJQS(output string) {
	_, err := exec.LookPath("jq")

	if err != nil {
		fmt.Println(output)
	} else {
		cmd := exec.Command("jq", ".")

		stdin, err := cmd.StdinPipe()
		if err != nil {
			panic(err)
		}
        cmd.Stdout = os.Stdout

        err = cmd.Start()
		if err != nil {
			panic(err)
		}
        stdin.Write([]byte(output))
        stdin.Close()

        if err != nil {
            panic(err)
        }
	}
}
