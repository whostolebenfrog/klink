package asgard

import (
	"fmt"
	common "nokia.com/klink/common"
	console "nokia.com/klink/console"
)

// Returns the asgard baseurl for the given environment, or poke if not known / blank
func baseUrl(env string) string {
	if env == "prod" {
		return "http://btmgsrvasgard02.brislabs.com:8080/eu-west-1/"
	}
	return "http://asgard.brislabs.com:8080/eu-west-1/"
}

// Returns the asgard loadbalancer url for the service environment pair
func loadBalancer(service string, env string) string {
	return baseUrl(env) + "loadBalancer/show/" + service + ".json"
}

// Exits the program with an error if the supplied info doesn't validate
func validateListServers(srv string, env string) {
	if srv == "" {
		console.Fail(`
    You must supply a service to list the servers for.
    e.g. "klink list-servers ditto"
        `)
	}
}

type LoadBalancerResponse struct {
    
}

// List server information for the supplied service and environment pair
func ListServers(args common.Command) {
	srv := args.SecondPos
	env := args.ThirdPos

	validateListServers(srv, env)

	console.Green()
	fmt.Print(fmt.Sprintf("\n\t%s", srv))
	console.Reset()
	fmt.Print(" : ")
	console.Red()
	fmt.Print(env + "\n\n")
	console.Reset()

	fmt.Println(common.GetString(loadBalancer(srv, env)))
}
