package gateway

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// JSONBody is a request body made of JSON
type JSONBody map[string]interface{}

// FormBody is a request body for a form
type FormBody map[string]interface{}

// IClient is the inteface of a HTTP client.
type IClient interface {
	POST(string, interface{}) (int, []byte, error)
	GET(string) (int, []byte, error)
}

/**
 * Client
 */

// NewClient creates a new auth'd HTTP client.
func NewClient(token string) Client {
	return Client{http.Client{}, token}
}

// Client is a simple helper for making auth'd HTTP json requests.
type Client struct {
	client http.Client
	token  string
}

// Encodes a FormBody into a URL string. Handles first layer
func (c Client) encodeForm(body FormBody) []byte {
	if len(body) == 0 {
		return []byte{}
	}

	var str string
	i := 0
	l := len(body)
	for k, v := range body {
		str += string(k)
		str += "="
		str += v.(string)
		if i != l-1 {
			str += "&"
		}
		i++
	}

	return []byte(str)
}

func (c Client) transformBody(body interface{}) ([]byte, error) {
	var bytes []byte
	var err error
	switch body.(type) {
	case JSONBody:
		bytes, err = json.Marshal(body)
		break

	case FormBody:
		bytes = c.encodeForm(body.(FormBody))
		break
	}

	return bytes, err
}

func (c Client) request(method string, url string, body interface{}) (int, []byte, error) {
	jsonBytes, err := c.transformBody(body)
	if err != nil {
		return 0, []byte{}, err
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return 0, []byte{}, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.token))

	switch body.(type) {
	case JSONBody:
		req.Header.Add("Content-Type", "application/json")
		break

	case FormBody:
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		break
	}

	response, err := c.client.Do(req)
	if err != nil {
		return 0, []byte{}, err
	}
	defer response.Body.Close()

	bod, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return 0, []byte{}, err
	}

	return response.StatusCode, bod, nil
}

// POST sends a POST request with the body. Returns status code and body.
func (c Client) POST(url string, body interface{}) (int, []byte, error) {
	return c.request("POST", url, body)
}

// GET sends a GET request. Returns status code and body.
func (c Client) GET(url string) (int, []byte, error) {
	return c.request("GET", url, map[string]interface{}{})
}
