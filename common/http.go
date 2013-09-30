package common

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

func PostJson(url string, body interface{}) (string, error) {
	b, err := json.Marshal(body)
	if err != nil {
		return "", errors.New("Unable to Marshall json for http post")
	}

	resp, err := http.Post(url, "application/json", bytes.NewReader(b))
	if err != nil {
		return "", errors.New(fmt.Sprintf("Error trying to call URL: %s", url))
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 || resp.StatusCode == 201 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", errors.New(fmt.Sprintf("Failed to read response body from: %s", url))
		}
		return string(body), nil
	} else {
		return "", errors.New(fmt.Sprintf("Got non 200 series response calling:", url, "with body", b))
	}
	return "", errors.New("I didn't think this was reachable :-(")
}

func GetString(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Error calling URL: %s", url))
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", errors.New(fmt.Sprintf("Failed to read body response from: %s", url))
		}
        return string(body), nil
	} else {
        return "", errors.New(fmt.Sprintf("Non 200 response from: %s", url))
    }
}

func GetJson(url string, v interface{}) error {
    body, err := GetString(url)
    if err != nil {
        return err
    }
    return json.Unmarshal([]byte(body), &v)
}
