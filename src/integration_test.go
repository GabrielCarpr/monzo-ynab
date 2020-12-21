package main_test

import (
	"monzo-ynab/internal/app"
	"monzo-ynab/internal/tester"
	"monzo-ynab/rest"
	"monzo-ynab/ynab"
	"testing"
)

func TestMonzoWebhook(t *testing.T) {
	mock := tester.ClientMock{}
	mock.PostStatus = 201

	builder := app.BuildApp(tester.GetTestConfig())
	tester.SetDep(builder, "ynab-client", &mock)
	handler := builder.Build().Get("rest-handler").(*rest.Handler)

	monzoEvent := map[string]interface{}{
		"type": "transaction.created",
		"data": map[string]interface{}{
			"id":          "123",
			"description": "test",
			"amount":      6900,
			"settled":     "2020-12-25T12:00:59Z",
		},
	}

	status, _ := tester.Request(handler, "POST", "/events/monzo/", monzoEvent)

	if status != 200 {
		t.Errorf("Bad status: %v", status)
	}

	received := mock.PostReceived["transaction"].(ynab.Transaction)

	if received.ImportID != "monzo:123" {
		t.Errorf("Wrong import ID: %s", received.ImportID)
	}
	if received.AccountID != tester.GetTestConfig().YNABAccountID {
		t.Errorf("Wrong account ID: %s", received.AccountID)
	}
	if received.PayeeName != "test" {
		t.Errorf("Wrong payee: %s", received.PayeeName)
	}
	if received.Date != "2020-12-25" {
		t.Errorf("Wrong date: %s", received.Date)
	}
	if received.Amount != 69000 {
		t.Errorf("Wrong amount: %v", received.Amount)
	}
	if received.Memo != "test" {
		t.Errorf("Wrong memo: %s", received.Memo)
	}
	if received.Cleared != "cleared" {
		t.Errorf("Wrong cleared value: %s", received.Cleared)
	}
	if received.Approved != false {
		t.Errorf("Wrong approved value: %v", received.Approved)
	}
}
