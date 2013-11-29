package console

import (
    "net/http"
    "net/url"
    common "nokia.com/klink/common"
)

func Hubot(message string, args common.Command) int {
    if args.Silent {
        return 0
    }
    resp, _ := http.PostForm("http://btmgsrvhubot001.brislabs.com/hubot/say",
         url.Values{"room" : {"503594"}, "message" : {message}})
    return resp.StatusCode
}

func Speak(args common.Command) {
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
    }
    http.PostForm("http://btmgsrvhubot001.brislabs.com/hubot/say",
         url.Values{"room" : {room}, "message" : {args.Message}})
}
