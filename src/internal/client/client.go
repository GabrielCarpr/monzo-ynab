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
	POST(string, map[string]interface{}) (int, map[string]interface{}, error)
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

func (c Client) request(method string, url string, body map[string]interface{}) (int, map[string]interface{}, error) {
	jsonBytes, err := json.Marshal(body)
	if err != nil {
		return 0, map[string]interface{}{}, err
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonBytes))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.token))
	if err != nil {
		return 0, map[string]interface{}{}, err
	}

	response, err := c.client.Do(req)
	defer response.Body.Close()
	if err != nil {
		return 0, map[string]interface{}{}, err
	}

	bod, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return 0, map[string]interface{}{}, err
	}

	var resBody map[string]interface{}
	if err := json.Unmarshal(bod, &resBody); err != nil {
		return 0, map[string]interface{}{}, err
	}

	return response.StatusCode, resBody, nil
}

// POST sends a POST request with the body. Returns status code and body.
func (c Client) POST(url string, body map[string]interface{}) (int, map[string]interface{}, error) {
	return c.request("POST", url, body)
}

// GET sends a GET request. Returns status code and body.
func (c Client) GET(url string) (int, map[string]interface{}, error) {
	return c.request("GET", url, map[string]interface{}{})
}
