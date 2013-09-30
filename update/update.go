package update

import (
    "bytes"
    "fmt"
    "os/exec"
    "runtime"
    common "nokia.com/klink/common"
    console "nokia.com/klink/console"
)

// TODO: version in a property?
// TODO: version bumped on build
// TODO: support non 0.x versions

const version = 4

func Version() {
    fmt.Println(fmt.Sprintf("klink-0.%d-%s-%s", version, runtime.GOOS, runtime.GOARCH))
}

func Update(path string) {
    currentVersion := version + 1

    nextVersion := fmt.Sprintf("klink-0.%d-%s-%s", currentVersion, runtime.GOOS, runtime.GOARCH)

    nextVersionUrl := "http://benkins.brislabs.com/klink/" + nextVersion

    exists, err := common.Head(nextVersionUrl)
    if err != nil {
        console.Fail("Failed to curl benkins for new version")
    }

    if exists {
        wget := exec.Command("wget", nextVersionUrl, "-O", path + ".tmp")
        var wgetStderr bytes.Buffer
        wget.Stderr = &wgetStderr
        wgetErr := wget.Run()

        if wgetErr != nil {
            fmt.Println(fmt.Sprint(wgetErr) + ":" + wgetStderr.String())
            console.Fail("Failed to wget the latest version. Ensure it's installed!")
        }

        mv := exec.Command("mv", "-f", path + ".tmp", path)
        var mvStderr bytes.Buffer
        mv.Stderr = &mvStderr
        mvErr := mv.Run()

        if mvErr != nil {
            fmt.Println(fmt.Sprint(mvErr) + ":" + mvStderr.String())
            console.Fail("Can't overwrite the previous version. DIY.")
        }

        chmod := exec.Command("chmod", "+x", path)
        var chmodStderr bytes.Buffer
        chmod.Stderr = &chmodStderr
        chmodErr := chmod.Run()

        if chmodErr != nil {
            fmt.Println(fmt.Sprint(chmodErr) + ":" + chmodStderr.String())
            console.Fail("Failed to set executable on the newly updated klink")
        }

        fmt.Println("Klink has been updated to the latest version!")
    } else {
        fmt.Println("You are using the latest version already. Good work kid, don't get cocky.")
    }
}
