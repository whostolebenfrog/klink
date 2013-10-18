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

// Perform an HTTP PUT on the supplied url with the body of the supplied object reference
// Returns a non nil error on non 200 series response or other error.
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

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Failed to read response body from: %s", url))
	}
	if resp.StatusCode == 200 || resp.StatusCode == 201 {
		return string(responseBody), nil
	}
	fmt.Println(string(responseBody))
	fmt.Println("%d response calling URL: ", resp.StatusCode, resp.StatusCode)
	return string(responseBody),
		errors.New(fmt.Sprintf("Got %d series response calling: %s with body: %s",
			resp.StatusCode, url, string(b)))
}

// Posts the supplied JSON to the url and unmarshals the response to the supplied
// struct.
func PostJsonUnmarshalResponse(url string, body interface{}, v interface{}) error {
	responseBody, err := PostJson(url, &body)
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(responseBody), &v)
}

// Performs an HTTP PUT on the supplied url with the body of the supplied object reference
// Returns a non nil error on non 200 series response or other error.
func PutJson(url string, body interface{}) (string, error) {
	b, err := json.Marshal(body)
	if err != nil {
		fmt.Println("Can't marshall body")
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

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Failed to read response body from: %s", url))
	}

	if resp.StatusCode == 200 || resp.StatusCode == 201 {
		return string(responseBody), nil
	}
	fmt.Println("%d response calling URL: %s", resp.StatusCode, url)
	return string(responseBody), errors.New(fmt.Sprintf("Got %d response calling: %s with body: %s",
		resp.StatusCode, url, b))
}

// Performs an HTTP GET request on the supplied url and returns the result
// as a string. Returns non nil err on non 200 response.
func GetString(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Error calling URL: %s", url))
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Failed to read body response from: %s", url))
	}

	if resp.StatusCode == 200 {
		return string(body), nil
	}
	fmt.Println(string(body))
	return string(body), errors.New(fmt.Sprintf("%d response from: %s", resp.StatusCode, url))
}

// Performs an HTTP GET request on the supplied url and unmarshals the response
// into the supplied object. Returns non nil error on failure
func GetJson(url string, v interface{}) error {
	body, err := GetString(url)
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(body), &v)
}

// Performs an HTTP HEAD call on the supplied URL. Returns true if
// the response code is 200.
func Head(url string) (bool, error) {
	resp, err := http.Head(url)
	if err != nil {
		return false, errors.New(fmt.Sprintf("Error calling head on URL: %s", url))
	}
	switch resp.StatusCode {
	case 200:
		return true, nil
	case 404:
		return false, nil
	default:
		panic(fmt.Sprintf("Unknown response: %d from HEAD on URL: %s Is your proxy set correctly?",
			resp.StatusCode, url))
	}

	return resp.StatusCode == 200, nil
}

// Config defines the configuration for a TimeoutDialer or TimeoutClient
type Config struct {
	ConnectTimeout   time.Duration
	ReadWriteTimeout time.Duration
}

// Creates a new TimeoutDialer with supplied config object
func TimeoutDialer(config *Config) func(net, addr string) (c net.Conn, err error) {
	return func(netw, addr string) (net.Conn, error) {
		conn, err := net.DialTimeout(netw, addr, config.ConnectTimeout)
		if err != nil {
			return nil, err
		}

		conn.SetDeadline(time.Now().Add(config.ReadWriteTimeout))
		conn.SetReadDeadline(time.Now().Add(config.ReadWriteTimeout))
		conn.SetWriteDeadline(time.Now().Add(config.ReadWriteTimeout))
		return conn, nil
	}
}

// Creates a new http.Client with a longer timeout, optionally accepts a Config object
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
