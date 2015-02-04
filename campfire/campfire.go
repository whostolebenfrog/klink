package campfire

import (
	"io/ioutil"
	"os"

	common "mixrad.io/klink/common"
	console "mixrad.io/klink/console"
)

func Init() {
	common.Register(
		common.Component{"campfire", Campfire,
			"{room} pipe stdin to a campfire room", "ROOMS"})
}

// writes stdin to the supplied campfire room
func Campfire(args common.Command) {
	room := args.SecondPos
	if room == "" {
		console.FailWithValidRooms(".. you didn't even pass a room! ..")
	}

	bytes, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	strBytes := string(bytes)
	if len(strBytes) == 0 {
		console.Fail(
			"You didn't pass any bytes. You're no Luis Suarez!*\n\n*(topical at time of writing)")
	}

	console.Speak(room, strBytes)
}
