package console

import (
	"fmt"
	"net/http"
	"net/url"

	common "nokia.com/klink/common"
)

func Init() {
	common.Register(common.Component{"speak", Speak, "Here be dragons"})
}

func doSpeak(room string, message string) int {
	resp, _ := http.PostForm("http://btmgsrvhubot001.brislabs.com/hubot/say",
		url.Values{"room": {room}, "message": {message}})
	return resp.StatusCode
}

func Hubot(message string, args common.Command) int {
	if args.Silent {
		return 0
	}
	return doSpeak("503594", message)
}

func Speak(args common.Command) {
	if args.Message == "" {
		fmt.Println("You fail. DRAGONS I SAID.")
	}
	room := ""
	switch args.SecondPos {
	case "general":
		room = "503594"
	case "clojure":
		room = "529176"
	case "cloud":
		room = "574028"
	case "asimov":
		room = "551265"
	case "kafka":
		room = "575611"
	case "hack":
		room = "582412"
	default:
		Fail("Unknown room")
	}
	doSpeak(room, args.Message)
}
