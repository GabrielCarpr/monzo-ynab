package gateway

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// IClient is the inteface of a HTTP client.
type IClient interface {
	POST(string, map[string]interface{}) (int, []byte, error)
	GET(string) (int, []byte, error)
}

// NewClient creates a new auth'd HTTP client.
func NewClient(token string) Client {
	return Client{http.Client{}, token}
}

// Client is a simple helper for making auth'd HTTP json requests.
type Client struct {
	client http.Client
	token  string
}

func (c Client) request(method string, url string, body map[string]interface{}) (int, []byte, error) {
	jsonBytes, err := json.Marshal(body)
	if err != nil {
		return 0, []byte{}, err
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return 0, []byte{}, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.token))

	response, err := c.client.Do(req)
	defer response.Body.Close()
	if err != nil {
		return 0, []byte{}, err
	}

	bod, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return 0, []byte{}, err
	}

	return response.StatusCode, bod, nil
}

// POST sends a POST request with the body. Returns status code and body.
func (c Client) POST(url string, body map[string]interface{}) (int, []byte, error) {
	return c.request("POST", url, body)
}

// GET sends a GET request. Returns status code and body.
func (c Client) GET(url string) (int, []byte, error) {
	return c.request("GET", url, map[string]interface{}{})
}
