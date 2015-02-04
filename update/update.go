package update

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strconv"
	"strings"
	"time"

	common "mixrad.io/klink/common"
	complete "mixrad.io/klink/complete"
	conf "mixrad.io/klink/conf"
	console "mixrad.io/klink/console"
	props "mixrad.io/klink/props"
)

func Init() {
	common.Register(
		common.Component{"update", Update,
			"Update klink to the latest version", ""},
		common.Component{"force-update", ForceUpdate,
			"Force klink to update to the current version", ""})
	/*
		common.Register(
			common.Component{"update", Update,
				"Update klink to the latest version"},
			common.Component{"force-update", ForceUpdate,
				"Force klink to update to the current version"})*/
}

func benkinsUrl(end string) string {
	return conf.DeployUrl + end
}

// Prints the current version, os and architecture
func PrintVersion() {
	fmt.Printf("klink-0.%d-%s-%s\n", Version, runtime.GOOS, runtime.GOARCH)
}

// Return the latest released version of klink
func LatestVersion() int {
	latestFromServer := common.GetString(benkinsUrl("version"))

	i, err := strconv.Atoi(strings.Replace(latestFromServer, "\n", "", 1))
	if err != nil {
		fmt.Println(err)
		console.Fail("Unable to get latest version. Check " + conf.DeployUrl)
	}
	return i
}

// Prints out an update error along with the path to manually get the
// next version
func errorWithHelper(nextVersionUrl string) {
	fmt.Println("\nThere appears to be a later version but an error has occured whilst updating")
	fmt.Println("You should may be able to download it manually from: ", nextVersionUrl)
	fmt.Println("\nYou could also try again or check benkins manually for updates")
	console.Fail(conf.DeployUrl)
}

// Update if there is a later version available, takes the path that this
// command was run from which is used as a backup if klink can't be
// found on the path
func Update(args common.Command) {
	complete.GenComplete(args)

	argsPath := os.Args[0]

	path, pathErr := exec.LookPath(path.Base(argsPath))
	if pathErr != nil {
		path = argsPath
	}

	props.SetLastUpdated(int32(time.Now().Unix()))

	latestVersion := LatestVersion()

	if latestVersion == Version {
		fmt.Println("You are using the latest version already. Good work kid, don't get cocky.")
		PrintVersion()
		return
	}

	nextVersion := fmt.Sprintf("klink-%d-%s-%s", latestVersion, runtime.GOOS, runtime.GOARCH)
	if common.IsWindows() {
		nextVersion += ".exe"
	}
	nextVersionUrl := benkinsUrl(nextVersion)

	if common.Head(nextVersionUrl) {
		doUpdate(nextVersionUrl, path, latestVersion)
	} else {
		errorWithHelper(nextVersionUrl)
	}
}

// For testing the update functionality
func ForceUpdate(_ common.Command) {
	argsPath := os.Args[0]

	path, pathErr := exec.LookPath(path.Base(argsPath))
	if pathErr != nil {
		path = argsPath
	}

	thisVersion := fmt.Sprintf("klink-%d-%s-%s", Version, runtime.GOOS, runtime.GOARCH)
	thisVersionUrl := benkinsUrl(thisVersion)

	if common.Head(thisVersionUrl) {
		doUpdate(thisVersionUrl, path, LatestVersion())
	} else {
		errorWithHelper(thisVersionUrl)
	}
}

// Does the update
func doUpdate(nextVersionUrl string, path string, latestVersion int) {
	resp, err := http.Get(nextVersionUrl)
	if err != nil {
		errorWithHelper(nextVersionUrl)
	}
	defer resp.Body.Close()

	file, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		errorWithHelper(nextVersionUrl)
	}

	err = ioutil.WriteFile(path+".update", file, 0755)
	if err != nil {
		fmt.Println(err)
		errorWithHelper(nextVersionUrl)
	}

	console.Green()
	fmt.Println(fmt.Sprintf("Klink has been updated to the latest version! %d to %d", Version, latestVersion))
	console.Reset()
	if common.IsWindows() {
		deferCopyForWindows(nextVersionUrl, path)
	} else {
		deferCopy(nextVersionUrl, path)
	}
}

// Write and run a script to copy the new version over ourselves, avoids
// file locks
func deferCopyForWindows(nextVersionUrl string, path string) {
	script := "Start-sleep 1\r\n" + "rm " + path + "\r\n" + "mv " + path + ".update " + path
	scriptBytes := []byte(script)
	ioutil.WriteFile("updateklink.PS1", scriptBytes, 0755)

	cmd := exec.Command("powershell", "-ExecutionPolicy", "ByPass", "-File", "updateklink.PS1")
	err := cmd.Start()
	if err != nil {
		fmt.Println(err)
		errorWithHelper(nextVersionUrl)
	}
	os.Exit(0)
}

// Write and run a script to copy the new version over ourselves, avoids
// file locks
func deferCopy(nextVersionUrl string, path string) {
	script := "sleep 1"
	script += "\nmv " + path + ".update " + path
	if len(os.Args) > 1 && os.Args[1] != "update" {
		script += "\n" + path + " "
		for i := range os.Args[1:] {
			argy := os.Args[i+1]
			if strings.Contains(argy, " ") {
				script += "\"" + argy + "\" "
			} else {
				script += argy + " "
			}
		}
	}
	script += "\nrm -f updateklink.sh"
	scriptBytes := []byte(script)
	ioutil.WriteFile("updateklink.sh", scriptBytes, 0755)

	cmd := exec.Command("sh", "updateklink.sh")
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}

	os.Exit(0)
}

// If we haven't tried to update recently then run an update first
func EnsureUpdatedRecently(argsPath string) {
	lastUpdated := props.GetLastUpdated()
	if lastUpdated == 0 {
		props.SetLastUpdated(int32(time.Now().Unix()))
		if LatestVersion() != Version {
			Update(common.Command{})
		}
	}

	now := int32(time.Now().Unix())
	if (now - lastUpdated) > (60 * 60 * 1) {
		if LatestVersion() != Version {
			Update(common.Command{})
		}
	}
}
