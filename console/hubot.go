package console

import (
    "net/http"
    "net/url"
)

func Hubot(message string) int {
    resp, _ := http.PostForm("http://btmgsrvhubot001.brislabs.com/hubot/say",
         url.Values{"room" : {"503594"}, "message" : {message}})
    return resp.StatusCode
}
