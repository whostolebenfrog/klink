package console

import (
	"fmt"
	"os"
)

func Klink() {
	fmt.Println("\n")
    fmt.Println("                        888      888 d8b          888")
	fmt.Println("   _   _                888      888 Y8P          888")
	fmt.Println("  ( \\_/ )  _   _        888      888              888")
	fmt.Println(" __) _ (__( \\_/ )       888  888 888 888 88888b.  888  888")
	fmt.Println("(__ (_) __)) _ (__      888 .88P 888 888 888 \"88b 888 .88P")
	fmt.Println("   ) _ ((__ (_) __)     888888K  888 888 888  888 888888K")
	fmt.Println("  (_/ \\_)  ) _ (        888 \"88b 888 888 888  888 888 \"88b")
	fmt.Println("          (_/ \\_)       888  888 888 888 888  888 888  888")
	fmt.Println("                         ...  ... ... ... ...  ... ...  ...")
	fmt.Println("\n")
}

func FailWhale(message string) {
	fmt.Println("\t     FAIL WHALE!")
	fmt.Println("\t")
	fmt.Println("\tW     W      W        ")
	fmt.Println("\tW        W  W     W    ")
	fmt.Println("\t              '.  W      ")
	fmt.Println("\t  .-\"\"-._     \\ \\.--|  ")
	fmt.Println("\t /       \"-..__) .-'   ")
	fmt.Println("\t|     _         /      ")
	fmt.Println("\t\\'-.__,   .__.,'       ")
	fmt.Println("\t `'----'._\\--'      ")
	fmt.Println("\tVVVVVVVVVVVVVVVVVVVVV")
	fmt.Println("\n")
	fmt.Println(message)
	fmt.Println("\n")
}

func Fail(message string) {
	fmt.Println(message)
	os.Exit(1)
}

func BigFail(message string) {
	FailWhale(message)
	os.Exit(1)
}
