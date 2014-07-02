package common

import (
	"fmt"
	"os"
	"regexp"
)

// Fail if the supplied string doesn't look like an ami
func FailIfNotAmi(maybeAmi string) {
	if maybeAmi != "" {
		matched, err := regexp.MatchString("^ami-.+$", maybeAmi)
		if err != nil {
			panic(err)
		}
		if !matched {
			fmt.Println(fmt.Sprintf("%s Doesn't look like an ami", maybeAmi))
			os.Exit(1)
		}
	}

}
