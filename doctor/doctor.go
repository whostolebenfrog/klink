package doctor

import (
	"fmt"
	"net/http"
	common "nokia.com/klink/common"
	console "nokia.com/klink/console"
	props "nokia.com/klink/props"
	"os"
	"strings"
	"time"
)

func Init() {
	common.Register(
		common.Component{"doctor", Doctor, "Test that everything is setup for klink to function"})
}

// Perform some doctor'in
func Doctor(args common.Command) {
	SetDoctorRun()

	console.Brown()
	fmt.Println("\n\t...The Doctor is in the house...\n")
	console.Reset()

	BrislabsReachable()

	console.Green()
	console.Bold()
	fmt.Println("\nThat's everything for now. You're good to go.")
	console.Reset()
}

// Checks to make sure that brislabs aws is reachable. Catches proxy issues mostly.
func BrislabsReachable() {
	// this handles the case where we fail to connect inside the timeout window
	defer func() {
		if p := recover(); p != nil {
			console.Red()
			fmt.Println("Failed")
			console.Reset()
			CheckProxySettings()
		}
	}()

	fmt.Println("Checking to see if aws resources in Brislabs are reachable...")

	httpClient := common.NewTimeoutClient(2*time.Second, 2*time.Second)
	req, err := http.NewRequest("HEAD", "http://ditto.brislabs.com:8080/1.x/status", nil)
	if err != nil {
		panic(err)
	}

	_, err = httpClient.Do(req)

	if err != nil {
		panic(err)
	}

	// if we made it here we could connect
	console.Green()
	fmt.Println("OK")
	console.Reset()
}

func EnsureDoctorRun() {
	if !HasDoctorRun() {
		SetDoctorRun()
	}
}

func SetDoctorRun() {
	props.Set("DoctorHasRun", "true")
}

func HasDoctorRun() bool {
	return props.GetRCProperties().DoctorHasRun == "true"
}

// Fail with doctor run message
func doctorOut() {
	console.Green()
	fmt.Println("\nRe-run the doctor with 'klink doctor'.")
	console.Reset()
	console.Fail("")
}

// failure to fix
func unableToFixConnection() {
	console.Red()
	fmt.Println("\nUnable to fix your connection issues. Try curling:")
	fmt.Println("http://ditto.brislabs.com:8080/1.x/status")
	fmt.Println("This can't go through a proxy. If it's trying - that's your problem.")
	fmt.Println("Usually the fix is putting export no_proxy=$no_proxy,brislabs.com")
	fmt.Println("into your .bashrc and either sourcing it or starting a new shell.")
	fmt.Println("It's also possible that ditto is down. If so that's your problem.")
	console.Reset()
	doctorOut()
}

// Can we tell if the proxy is fucked?
func CheckProxySettings() {
	// can't do much for windows.
	if common.IsWindows() {
		fmt.Println("\nUnable to talk to brislabs.com in aws and unable to automatically fix it.")
		fmt.Println("It's likely that you have a proxy issue. Klink needs to connect directly")
		fmt.Println("to aws brislabs without a proxy. There might also be windows firewall")
		fmt.Println("issues. Talk to Ben Griffiths if you have these issues.")
		doctorOut()
	}

	proxy := os.Getenv("http_proxy") + os.Getenv("HTTP_PROXY")

	if proxy != "" {
		fmt.Printf("\nLooks like you have a proxy set: %s\n", proxy)

		noProxy := os.Getenv("no_proxy") + os.Getenv("NO_PROXY")
		if !strings.Contains(noProxy, "brislabs") {
			fmt.Printf("\nYou don't have brislabs in your no proxy list: %s\n", noProxy)

			BrislabsProxyWrong()
		} else {
			unableToFixConnection()
		}
	}
	unableToFixConnection()
}

func BrislabsProxyWrong() {
	console.Green()
	rc_message := `
Do you want me to add brislabs.com to no_proxy in your .bashrc?
No guaranty, if it breaks your bash file it's not my fault!
Please type Yes or No:`
	fmt.Println(rc_message)
	console.Reset()

	var response string

	fmt.Scan(&response)

	switch response {
	case "yes", "Yes", "YES", "y", "Y":
		SetBrislabsNoProxy()
	case "no", "No", "NO", "n", "N":
		unableToFixConnection()
	default:
		fmt.Println("Type better.")
		BrislabsProxyWrong()
	}
}

func SetBrislabsNoProxy() {
	bashRc := common.UserHomeDir() + "/.bashrc"

	f, err := os.OpenFile(bashRc, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	no_proxy := `
# Added by klink
# Setting no proxy for brislabs.com as requests to brislabs in aws will not resolve
# via a proxy.
export no_proxy=$no_proxy,brislabs.com`
	if _, err = f.WriteString(no_proxy); err != nil {
		panic(err)
	}

	console.Green()
	fmt.Println("\nOK - everything should be in order")
	fmt.Println("Source you bash.rc 'source ~/.bashrc'")
	fmt.Println("Then retry klink doctor")
	console.Reset()

	os.Exit(0)
}
