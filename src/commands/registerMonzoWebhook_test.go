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
	register := builder.Build().Get("register-webhook-command").(*commands.RegisterMonzoWebhook)

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

func Test_RegisterWebhookIsIdempotent(t *testing.T) {
	mock := clientMock{}
	mock.getStatus = 200
	mock.getResponse = map[string]interface{}{
		"webhooks": []monzo.Webhook{
			{
				AccountID: "test_monzo_id",
				URL:       "http://testurl/test",
				ID:        "test",
			},
			{
				AccountID: "test_monzo_id",
				URL:       "http://testurl/hello",
				ID:        "lol",
			},
		},
	}

	builder := app.BuildApp(getTestConfig())
	tester.SetDep(builder, "monzo-client", &mock)
	cmd := builder.Build().Get("register-webhook-command").(*commands.RegisterMonzoWebhook)

	err := cmd.Execute("/test")
	if err != nil {
		t.Errorf("Returned error: %s", err)
	}

	if len(mock.postReceived) > 0 {
		t.Errorf("Sent a post request: %v", mock.postReceived)
	}
}
