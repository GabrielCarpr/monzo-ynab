package commands_test

import (
	"monzo-ynab/commands"
	"monzo-ynab/internal/app"
	"monzo-ynab/internal/tester"
	"testing"

	"github.com/mitchellh/mapstructure"
)

func Test_ConvertMovesATransaction(t *testing.T) {
	mock := clientMock{}
	mock.getStatus = 200
	mock.getResponse = map[string]interface{}{
		"transaction": map[string]interface{}{
			"id":          "abc",
			"amount":      500,
			"settled":     "2019-12-25T12:35:00Z",
			"description": "test transaction",
		},
	}
	mock.postStatus = 201

	builder := app.BuildApp(getTestConfig())
	tester.SetDep(builder, "monzo-client", &mock)
	tester.SetDep(builder, "ynab-client", &mock)
	cmd := builder.Build().Get("convert-transaction-command").(*commands.Convert)

	err := cmd.Execute("abc")
	if err != nil {
		t.Errorf("Error'd: %s", err)
	}

	var received map[string]interface{}
	mapstructure.Decode(mock.postReceived["transaction"], &received)
	if received["Amount"].(int) != 5000 {
		t.Errorf("Wrong amount: %v", mock.postReceived["Amount"])
	}
	if received["Memo"].(string) != "test transaction" {
		t.Errorf("Wrong memo: %s", mock.postReceived["Memo"])
	}
	if received["Date"].(string) != "2019-12-25" {
		t.Errorf("Wrong date: %s", mock.postReceived["Date"])
	}
	if received["ImportID"].(string) != "monzo:abc" {
		t.Errorf("Wrong ImportID: %s", mock.postReceived["ImportID"])
	}
}
