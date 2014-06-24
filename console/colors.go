package console

import (
	"fmt"

	common "nokia.com/klink/common"
)

// Sets color and weight back to defaults
func Reset() {
	if !common.IsWindows() {
		print("\033[0m")
	}
}

// Sets color and weight back to defaults and flushes
func FReset() {
	Reset()
	fmt.Println("")
}

// Sets color red
func Red() {
	if !common.IsWindows() {
		print("\033[31m")
	}
}

// Sets weight to bold
func Bold() {
	if !common.IsWindows() {
		print("\033[1m")
	}
}

// Sets color green
func Green() {
	if !common.IsWindows() {
		print("\033[32m")
	}
}

// Sets color brown
func Brown() {
	if !common.IsWindows() {
		print("\033[33m")
	}
}

// Sets color blue
func Blue() {
	if !common.IsWindows() {
		print("\033[34m")
	}
}

// Sets color magenta
func Magenta() {
	if !common.IsWindows() {
		print("\033[35m")
	}
}

// Sets color cyan
func Cyan() {
	if !common.IsWindows() {
		print("\033[36m")
	}
}

// Sets color yellow
func Yellow() {
	if !common.IsWindows() {
		print("\033[33m")
	}
}

// Sets color grey
func Grey() {
	if !common.IsWindows() {
		print("\033[37m")
	}
}
