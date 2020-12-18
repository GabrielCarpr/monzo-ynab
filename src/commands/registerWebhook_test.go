package commands_test

import (
	"monzo-ynab/commands"
	"monzo-ynab/internal/app"
	"monzo-ynab/internal/tester"
	"monzo-ynab/monzo"
	"testing"
)

func Test_RegistersWebhook(t *testing.T) {
	mock := clientMock{}
	mock.getStatus = 200
	mock.postStatus = 200
	mock.getResponse = map[string]interface{}{
		"webhooks": []monzo.Webhook{},
	}

	builder := app.BuildApp(getTestConfig())
	tester.SetDep(builder, "monzo-client", &mock)
	register := builder.Build().Get("register-webhook-command").(*commands.RegisterWebhook)

	err := register.Execute("/testing")

	if err != nil {
		t.Errorf("Error'd: %s", err)
	}

	if mock.postReceived["url"].(string) != "http://testurl/testing" {
		t.Errorf("Webhook URL wrong: %s", mock.postReceived["url"])
	}
	if mock.postReceived["account_id"].(string) != "test_monzo_id" {
		t.Errorf("Wrong account ID: %s", mock.postReceived["account_id"])
	}
}
