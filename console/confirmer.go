package console

import (
	"fmt"
)

type ColourFunc func()

func Confirmer(colourer ColourFunc, message string) {
	colourer()
	fmt.Println(message)
	Reset()

	var response string

	fmt.Scan(&response)

	switch response {
	case "yes", "Yes", "YES", "y", "Y":
		break
	case "no", "No", "NO", "n", "N":
		Red()
		Fail("Aborted.")
		Reset()
	default:
		fmt.Println("Type better.")
		Confirmer(colourer, message)
	}
}
