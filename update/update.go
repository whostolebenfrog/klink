package update

import (
	//"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	common "nokia.com/klink/common"
	console "nokia.com/klink/console"
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
    fmt.Println("in update")
	path, pathErr := exec.LookPath("klink")
    fmt.Println(argsPath)
	if pathErr != nil {
        fmt.Println("using args path")
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
        fmt.Println("about to update")
        fmt.Println(path)
		doUpdate(nextVersionUrl, path)
	} else {
		fmt.Println(err)
		errorWithHelper(nextVersionUrl)
	}
}

func doUpdate(nextVersionUrl string, path string) {
	// get the latest version, save to a tmp file
    fmt.Println(nextVersionUrl)

	resp, err := http.Get(nextVersionUrl)
	if err != nil {
		errorWithHelper(nextVersionUrl)
	}
	defer resp.Body.Close()
    fmt.Println("got the old body")

    file, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        errorWithHelper(nextVersionUrl)
    }

    fmt.Println("about to write!")
	err = ioutil.WriteFile(path, file, 0755)
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println("should have written")

    err = ioutil.WriteFile("C:\\Users\\bgriffit\\klink\\yolo.exe", file, 0755)
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println("le fin")

	/*
	   wget := exec.Command("wget", nextVersionUrl, "-O", path+".tmp")
	   var wgetStderr bytes.Buffer
	   wget.Stderr = &wgetStderr
	   wgetErr := wget.Run()

	   if wgetErr != nil {
	       fmt.Println(fmt.Sprint(wgetErr) + ":" + wgetStderr.String())
	       fmt.Println("Failed to wget the latest version. Ensure it's installed!")
	       errorWithHelper(nextVersionUrl)
	   }

	   // TODO: don't -f on windows
	   // overwrite the old version with the new one
	   mv := exec.Command("mv", "-f", path+".tmp", path)
	   var mvStderr bytes.Buffer
	   mv.Stderr = &mvStderr
	   mvErr := mv.Run()

	   if mvErr != nil {
	       fmt.Println(fmt.Sprint(mvErr) + ":" + mvStderr.String())
	       fmt.Println("Can't overwrite the previous version. You might be able to do it yourself")
	       errorWithHelper(nextVersionUrl)
	   }

	   // Don't do this on windows!
	   // make the new one executable
	   chmod := exec.Command("chmod", "+x", path)
	   var chmodStderr bytes.Buffer
	   chmod.Stderr = &chmodStderr
	   chmodErr := chmod.Run()

	   if chmodErr != nil {
	       fmt.Println(fmt.Sprint(chmodErr) + ":" + chmodStderr.String())
	       fmt.Println("Failed to +x on klink. You might be able to do it yourself")
	       errorWithHelper(nextVersionUrl)
	   }

	   fmt.Println("Klink has been updated to the latest version!") */
}
