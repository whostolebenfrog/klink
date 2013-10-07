package common

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

func PostJson(url string, body interface{}) (string, error) {

	b, err := json.Marshal(body)
	if err != nil {
		fmt.Println("Can't marshall body")
		return "", errors.New("Unable to Marshall json for http post")
	}

	resp, err := http.Post(url, "application/json", bytes.NewReader(b))
	if err != nil {
		fmt.Println("Error posting to url:", url)
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
	fmt.Println("Non 200 response calling URL: ", resp.StatusCode)
	return "", errors.New(fmt.Sprintf("Got non 200 series response calling:", url, "with body", b))
}

func PutJson(url string, body interface{}) (string, error) {
	b, err := json.Marshal(body)
	if err != nil {
		return "", errors.New("Unable to Marshall json for http put")
	}

	req, _ := http.NewRequest("PUT", url, bytes.NewReader(b))

	req.Header.Add("Content-Type", "application/json")

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

type Config struct {
	ConnectTimeout   time.Duration
	ReadWriteTimeout time.Duration
}

func TimeoutDialer(config *Config) func(net, addr string) (c net.Conn, err error) {
	return func(netw, addr string) (net.Conn, error) {
		conn, err := net.DialTimeout(netw, addr, config.ConnectTimeout)
		if err != nil {
			return nil, err
		}

		conn.SetDeadline(time.Time{})
        conn.SetReadDeadline(time.Time{})
        conn.SetWriteDeadline(time.Time{})
		return conn, nil
	}
}

func NewTimeoutClient(args ...interface{}) *http.Client {
	// Default configuration, lots of time for a streaming response to return
    // Needs as long as our longest bake
	config := &Config{
		ConnectTimeout:   5 * time.Second,
		ReadWriteTimeout: 1200 * time.Second,
	}

	// merge the default with user input if there is one
	if len(args) == 1 {
		timeout := args[0].(time.Duration)
		config.ConnectTimeout = timeout
		config.ReadWriteTimeout = timeout
	}

	if len(args) == 2 {
		config.ConnectTimeout = args[0].(time.Duration)
		config.ReadWriteTimeout = args[1].(time.Duration)
	}

	return &http.Client{
		Transport: &http.Transport{
			Dial: TimeoutDialer(config),
		},
	}
}
