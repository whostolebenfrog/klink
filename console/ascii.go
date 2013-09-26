package console

import (
	"fmt"
	"os"
)

func Klink() {
	fmt.Println("   _   _           ")
	fmt.Println("  ( \\_/ )  _   _   ")
	fmt.Println(" __) _ (__( \\_/ )       ____ ____ ____ ____ ____ ")
	fmt.Println("(__ (_) __)) _ (__     ||k |||l |||i |||n |||k ||")
	fmt.Println("   ) _ ((__ (_) __)    ||__|||__|||__|||__|||__||")
	fmt.Println("  (_/ \\_)  ) _ (       |/__\\|/__\\|/__\\|/__\\|/__\\|")
	fmt.Println("          (_/ \\_)  ")
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
