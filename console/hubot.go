package console

import (
	"fmt"
	"net/http"
	"net/url"
)

func doSpeak(room string, message string) int {
	resp, _ := http.PostForm("http://btmgsrvhubot001.brislabs.com/hubot/say",
		url.Values{"room": {room}, "message": {message}})
	return resp.StatusCode
}

var rooms = map[string]string{
	"general":   "503594",
	"clojure":   "529176",
	"cloud":     "574028",
	"asimov":    "551265",
	"fusion":    "551264",
	"reportlog": "594551",
	"kafka":     "575611",
	"hack":      "582412",
	"testing":   "568845",
	"github":    "597627"}

// speak in the supplied room with the supplied message
func Speak(room string, message string) {

	roomNumber := rooms[room]
	if roomNumber == "" {
		FailWithValidRooms(room)
	}
	doSpeak(roomNumber, message)
}

// fail and list valid rooms
func FailWithValidRooms(room string) {
	fmt.Println(fmt.Sprintf("Room: %s is not known\n", room))
	allRooms := make([]string, len(rooms))

	i := 0
	for key, _ := range rooms {
		allRooms[i] = key
		i++
	}
	Fail(fmt.Sprintf("Known rooms are: %s", allRooms))
}
