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
