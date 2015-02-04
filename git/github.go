package git

import (
	"fmt"
	"io/ioutil"
	"os"

	common "mixrad.io/klink/common"
	conf "mixrad.io/klink/conf"
	console "mixrad.io/klink/console"
)

type GistJson struct {
	Content     string `json:"content"`
	Description string `json:"description"`
	Name        string `json:"name"`
}

func githubTeamsUrl(end string) string {
	return conf.GithubTeamsUrl + end
}

// send stdin to a github gist
func Gist(args common.Command) {
	fileName := args.SecondPos
	description := args.ThirdPos

	if fileName == "" {
		console.Fail("You must pass a filename, use the extension to set the type")
	}

	bytes, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	createReq := GistJson{string(bytes), description, fileName}
	fmt.Println(common.PostJson(githubTeamsUrl("/create-gist"), createReq))
}
