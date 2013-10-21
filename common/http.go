package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

// Perform an HTTP PUT on the supplied url with the body of the supplied object reference
func PostJson(url string, body interface{}) string {
	b, err := json.Marshal(body)
	if err != nil {
        panic(fmt.Sprintf("Can't marshall body attempting to call %s", url))
	}

	resp, err := http.Post(url, "application/json", bytes.NewReader(b))
	if err != nil {
        panic(fmt.Sprintf("Error trying to call URL: %s", url))
	}
	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
        panic(fmt.Sprintf("Failed to read response body from: %s", url))
	}
	if resp.StatusCode == 200 || resp.StatusCode == 201 {
		return string(responseBody)
	}
	fmt.Println(string(responseBody))
    panic(fmt.Sprintf("%d response calling URL: ", resp.StatusCode, resp.StatusCode))
}

// Posts the supplied JSON to the url and unmarshals the response to the supplied
// struct.
func PostJsonUnmarshalResponse(url string, body interface{}, v interface{}) {
	responseBody := PostJson(url, &body)
    err := json.Unmarshal([]byte(responseBody), &v)
    if err != nil {
        panic(fmt.Sprintf("Unable to unmarshal response from %s", url))
    }
}

func PutByteArray(url string, data []byte) string {
	req, err := http.NewRequest("PUT", url, bytes.NewReader(data))
    if err != nil {
        panic(fmt.Sprintf("Error making PUT request to url: %s", url))
    }
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
        panic(fmt.Sprintf("Error trying to call URL: %s", url))
	}
	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
        panic(fmt.Sprintf("Failed to read response body from: %s", url))
	}

	if resp.StatusCode == 200 || resp.StatusCode == 201 {
		return string(responseBody)
	}
    panic(fmt.Sprintf("Got %d response calling: %s with body: %s", resp.StatusCode, url, data))
}

// Performs an HTTP PUT on the supplied url with the body of the supplied string.
func PutString(url string, body string) string {
    return PutByteArray(url, []byte(body))
}

// Performs an HTTP PUT on the supplied url with the body of the supplied object reference
func PutJson(url string, body interface{}) string {
	b, err := json.Marshal(body)
	if err != nil {
        panic(fmt.Sprintf("Unable to Marshall json for http put to url %s", url))
	}
    return PutByteArray(url, b)
}

// Performs an HTTP GET request on the supplied url and returns the result
// as a string. Returns non nil err on non 200 response.
func GetString(url string) string {
	resp, err := http.Get(url)
	if err != nil {
        panic(fmt.Sprintf("Error calling URL: %s", url))
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
        panic(fmt.Sprintf("Failed to read body response from: %s", url))
	}

	if resp.StatusCode == 200 {
		return string(body)
	}
    panic(fmt.Sprintf("Got %d response calling: %s. Response was:\n%s",
        resp.StatusCode, url, string(body)))
}

// Performs an HTTP GET request on the supplied url and unmarshals the response
// into the supplied object. Returns non nil error on failure
func GetJson(url string, v interface{}) {
	err := json.Unmarshal([]byte(GetString(url)), &v)
    if err != nil {
        panic(fmt.Sprintf("Umable to marshall response from url %s", url))
    }
}

// Performs an HTTP HEAD call on the supplied URL. Returns true if
// the response code is 200.
func Head(url string) bool {
	resp, err := http.Head(url)
	if err != nil {
        panic(fmt.Sprintf("Error calling head on URL: %s", url))
	}
	switch resp.StatusCode {
	case 200:
		return true
	case 404:
		return false
	default:
		panic(fmt.Sprintf("Unknown response: %d from HEAD on URL: %s Is your proxy set correctly?",
			resp.StatusCode, url))
	}
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
