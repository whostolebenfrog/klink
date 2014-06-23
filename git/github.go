package git

import (
	"fmt"
	"io/ioutil"
	common "nokia.com/klink/common"
	"os"
)

type File struct {
	Content string `json:"content"`
}

type GistJson struct {
	Description string `json:"description"`
	Public      string `json:"public"`
	Files       []File `json:"files"`
}

// send stdin to a github gist
func Gist(args common.Command) {
	bytes, err := ioutil.ReadAll(os.Stdin)
	fmt.Println(string(bytes))
	if err != nil {
		panic(err)
	}
}
