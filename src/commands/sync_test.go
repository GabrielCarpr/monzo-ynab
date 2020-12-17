package commands_test

import (
	"encoding/json"
	"monzo-ynab/commands"
	"monzo-ynab/internal/app"
	"monzo-ynab/internal/config"
	"monzo-ynab/internal/tester"
	"testing"

	"github.com/mitchellh/mapstructure"
)

func getTestConfig() config.Config {
	return config.Config{
		YNABToken:        "test",
		YNABAccountID:    "test_acc_id",
		YNABBudgetID:     "test_budget_id",
		MonzoAccountID:   "test_monzo_id",
		MonzoAccessToken: "test_monzo_token",
	}
}

type clientMock struct {
	postResponse map[string]interface{}
	postStatus   int
	postReceived map[string]interface{}

	getResponse map[string]interface{}
	getStatus   int
}

func (c *clientMock) POST(url string, body map[string]interface{}) (int, []byte, error) {
	c.postReceived = body
	res, _ := json.Marshal(c.postResponse)
	return c.postStatus, res, nil
}

func (c *clientMock) GET(url string) (int, []byte, error) {
	res, _ := json.Marshal(c.getResponse)
	return c.getStatus, res, nil
}

func Test_ImportsTransactions(t *testing.T) {
	mock := clientMock{}
	mock.postStatus = 201
	mock.getStatus = 200

	mock.getResponse = map[string]interface{}{
		"transactions": []map[string]interface{}{
			{
				"id":          "tx_ffi39fjkls",
				"description": "test transaction",
				"amount":      2500,
				"settled":     "2020-12-25T12:25:00Z",
			},
		},
	}

	builder := app.BuildApp(getTestConfig())
	tester.SetDep(builder, "ynab-client", &mock)
	tester.SetDep(builder, "monzo-client", &mock)
	ctn := builder.Build()
	sync := ctn.Get("sync-command").(*commands.Sync)

	err := sync.Execute(3)
	if err != nil {
		t.Errorf("Produced error: %s", err)
	}

	postReceived := mock.postReceived["transaction"]
	var received map[string]interface{}
	mapstructure.Decode(postReceived, &received)
	if received["Amount"].(int) != 25000 {
		t.Errorf("Wrong amount: %v", received["Amount"])
	}
	if received["Memo"].(string) != "test transaction" {
		t.Errorf("Wrong memo: %s", received["Memo"])
	}
	if received["Date"].(string) != "2020-12-25" {
		t.Errorf("Wrong date: %s", received["Date"])
	}
	if received["ImportID"].(string) != "monzo:tx_ffi39fjkls" {
		t.Errorf("ImportID wrong: %s", received["ImportID"])
	}
}
