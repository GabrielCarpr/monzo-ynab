package commands_test

import (
	"monzo-ynab/commands"
	"monzo-ynab/internal/app"
	"monzo-ynab/internal/tester"
	"testing"

	"github.com/mitchellh/mapstructure"
)

func Test_ImportsTransactions(t *testing.T) {
	mock := tester.ClientMock{}
	mock.PostStatus = 201
	mock.GetStatus = 200

	mock.GetResponse = map[string]interface{}{
		"transactions": []map[string]interface{}{
			{
				"id":          "tx_ffi39fjkls",
				"description": "test transaction",
				"amount":      2500,
				"settled":     "2020-12-25T12:25:00Z",
			},
		},
	}

	builder := app.BuildApp(tester.GetTestConfig())
	tester.SetDep(builder, "ynab-client", &mock)
	tester.SetDep(builder, "monzo-client", &mock)
	sync := builder.Build().Get("sync-command").(*commands.Sync)

	err := sync.Execute(3)
	if err != nil {
		t.Errorf("Produced error: %s", err)
	}

	postReceived := mock.PostReceived["transaction"]
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
