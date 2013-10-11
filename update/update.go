package update

import (
	"fmt"
	"io/ioutil"
	"net/http"
	common "nokia.com/klink/common"
	console "nokia.com/klink/console"
    "os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

func benkinsUrl(end string) string {
	return "http://benkins.brislabs.com/klink/" + end
}

// Prints the current version, os and architecture
func PrintVersion() {
	fmt.Println(fmt.Sprintf("klink-0.%d-%s-%s", Version, runtime.GOOS, runtime.GOARCH))
}

// Return the latest released version of klink
func LatestVersion() int {
	latestFromServer, err := common.GetString(benkinsUrl("version"))
	if err != nil {
		fmt.Println(err)
		console.Fail("Unable to get latest version. Check http://benkins.brislabs.com/klink/")
	}

	i, err := strconv.Atoi(strings.Replace(latestFromServer, "\n", "", 1))
	if err != nil {
		fmt.Println(err)
		console.Fail("Unable to get latest version. Check http://benkins.brislabs.com/klink/")
	}
	return i
}

// Prints out an update error along with the path to manually get the
// next version
func errorWithHelper(nextVersionUrl string) {
	fmt.Println("\nThere appears to be a later version but an error has occured whilst updating")
	fmt.Println("You should may be able to download it manually from: ", nextVersionUrl)
	fmt.Println("\nYou could also try again or check benkins manually for updates")
	console.Fail("http://benkins.brislabs.com/klink/")
}

// Update if there is a later version available, takes the path that this
// command was run from which is used as a backup if klink can't be
// found on the path
func Update(argsPath string) {
	path, pathErr := exec.LookPath("klink")
	if pathErr != nil {
		path = argsPath
	}

	latestVersion := LatestVersion()

    /*
	if latestVersion == Version {
		fmt.Println("You are using the latest version already. Good work kid, don't get cocky.")
		PrintVersion()
		return
	}*/

	nextVersion := fmt.Sprintf("klink-%d-%s-%s", latestVersion, runtime.GOOS, runtime.GOARCH)
	nextVersionUrl := benkinsUrl(nextVersion)

	exists, err := common.Head(nextVersionUrl)
	if err != nil {
		fmt.Println(err)
		errorWithHelper(nextVersionUrl)
	}

	if exists {
		doUpdate(nextVersionUrl, path)
	} else {
		fmt.Println(err)
		errorWithHelper(nextVersionUrl)
	}
}

func doUpdate(nextVersionUrl string, path string) {
	resp, err := http.Get(nextVersionUrl)
	if err != nil {
		errorWithHelper(nextVersionUrl)
	}
	defer resp.Body.Close()

    file, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        errorWithHelper(nextVersionUrl)
    }

    if !common.IsWindows() {
        err = ioutil.WriteFile(path, file, 0755)
        if err != nil {
            fmt.Println("here2")
            fmt.Println(err)
            errorWithHelper(nextVersionUrl)
        }
    } else {
        err = ioutil.WriteFile(path+".update", file, 0755)
        if err != nil {
            fmt.Println(err)
            errorWithHelper(nextVersionUrl)
        }
    }

    fmt.Println("Klink has been updated to the latest version!")
    if common.IsWindows() {
        deferCopyForWindows(nextVersionUrl, path)
    }
}

// Windows refuses to us overwrite ourselves so make a script with a small sleep
// that overwrites us with the new version. NEVER GIVE UP.
func deferCopyForWindows (nextVersionUrl string, path string) {
    script := "Start-sleep 1\n\r" + "mv " + path + ".update " + path
    scriptBytes := []byte(script)
    ioutil.WriteFile("update.PS1", scriptBytes, 0755)

    cmd := exec.Command("powershell", "-ExecutionPolicy", "ByPass", "-File", "update.PS1")
    err := cmd.Start()
    if err != nil {
        errorWithHelper(nextVersionUrl)
    }
    os.Exit(0)
}
