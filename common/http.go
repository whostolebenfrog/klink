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
	}
	return "", errors.New(fmt.Sprintf("Got non 200 series response calling:", url, "with body", b))
}

func PutJson(url string, body interface{}) (string, error) {
	b, err := json.Marshal(body)
	if err != nil {
		return "", errors.New("Unable to Marshall json for http put")
	}

	req, _ := http.NewRequest("PUT", url, bytes.NewReader(b))

	client := &http.Client{}

	resp, err := client.Do(req)

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
	}
	return "", errors.New(fmt.Sprintf("Got non 200 series response calling:", url, "with body", b))
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
	}
	return "", errors.New(fmt.Sprintf("Non 200 response from: %s", url))
}

func GetJson(url string, v interface{}) error {
	body, err := GetString(url)
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(body), &v)
}

func Head(url string) (bool, error) {
	resp, err := http.Head(url)
	if err != nil {
		return false, errors.New(fmt.Sprintf("Error calling head on URL: %s", url))
	}

	return resp.StatusCode == 200, nil
}
