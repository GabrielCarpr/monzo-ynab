package tester

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	client "monzo-ynab/internal/client"
	"monzo-ynab/internal/config"
	"net/http"
	"net/http/httptest"

	"github.com/sarulabs/di/v2"
)

// SetDep is a shorthand for setting a dependency on a DI builder
func SetDep(builder *di.Builder, name string, object interface{}) {
	builder.Add(di.Def{
		Name: name,
		Build: func(ctn di.Container) (interface{}, error) {
			return object, nil
		},
	})
}

// GetTestConfig returns a Config set up for testing.
func GetTestConfig() config.Config {
	return config.Config{
		YNABToken:        "test",
		YNABAccountID:    "test_acc_id",
		YNABBudgetID:     "test_budget_id",
		MonzoAccountID:   "test_monzo_id",
		MonzoAccessToken: "test_monzo_token",
		BaseURL:          "http://testurl",
	}
}

// ClientMock is a mock for mocking the internal/client.
type ClientMock struct {
	PostResponse map[string]interface{}
	PostStatus   int
	PostReceived map[string]interface{}

	GetResponse map[string]interface{}
	GetStatus   int
}

// POST mocks a POST request.
func (c *ClientMock) POST(url string, body interface{}) (int, []byte, error) {
	switch body.(type) {
	case client.JSONBody:
		c.PostReceived = body.(client.JSONBody)
		break
	case client.FormBody:
		c.PostReceived = body.(client.FormBody)
		break
	case map[string]interface{}:
		c.PostReceived = body.(map[string]interface{})
		break
	}

	res, _ := json.Marshal(c.PostResponse)
	return c.PostStatus, res, nil
}

// GET mocks a GET request.
func (c *ClientMock) GET(url string) (int, []byte, error) {
	res, _ := json.Marshal(c.GetResponse)
	return c.GetStatus, res, nil
}

// Request sends a test HTTP request to a provided handler.
func Request(handler http.Handler, method string, path string, jsonBody map[string]interface{}) (int, []byte) {
	json, err := json.Marshal(jsonBody)
	if err != nil {
		panic(err)
	}

	req := httptest.NewRequest(method, path, bytes.NewBuffer(json))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	response := w.Result()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	return response.StatusCode, body
}
