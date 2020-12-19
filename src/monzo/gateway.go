package monzo

import (
	"encoding/json"
	"fmt"
	"monzo-ynab/domain"
	client "monzo-ynab/internal/client"
	"monzo-ynab/internal/config"
	"time"
)

const monzoAPIUrl = "https://api.monzo.com"

// Transaction represents a Monzo API DTO (or, just the fields we need.)
type Transaction struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	Amount      int    `json:"amount"`
	Settled     string `json:"settled"`
}

type transactionsResponse struct {
	Transactions []Transaction `json:"transactions"`
}

type transactionResponse struct {
	Transaction Transaction `json:"transaction"`
}

// Transaction converts a Monzo transaction to a domain Transaction.
func (t Transaction) Transaction() domain.Transaction {
	date, _ := time.Parse(time.RFC3339, t.Settled)
	return domain.Transaction{
		Amount:      t.Amount,
		Date:        date,
		Payee:       t.Description,
		Description: t.Description,
		MonzoID:     t.ID,
	}
}

// NewGateway constructs a Monzo gateway.
func NewGateway(config config.Config, c client.IClient) *Gateway {
	return &Gateway{config, c}
}

// Gateway is a Monzo gateway.
type Gateway struct {
	config config.Config
	client client.IClient
}

// ListTransactions gets transactions since days.
func (g Gateway) ListTransactions(since string) ([]Transaction, error) {
	status, body, err := g.client.GET(fmt.Sprintf("%s/transactions?account_id=%s&since=%s", monzoAPIUrl, g.config.MonzoAccountID, since))
	if err != nil {
		return []Transaction{}, err
	}

	if status != 200 {
		return []Transaction{}, fmt.Errorf("API did not return 200: %v %v", status, string(body))
	}

	var tr transactionsResponse
	if err := json.Unmarshal(body, &tr); err != nil {
		return tr.Transactions, fmt.Errorf("Failed reading Monzo API response: %w", err)
	}

	return tr.Transactions, nil
}

// RegisterWebhook creates a webhook in Monzo.
func (g Gateway) RegisterWebhook(url string) error {
	status, body, err := g.client.POST(fmt.Sprintf("%s/webhooks", monzoAPIUrl), client.FormBody{
		"account_id": g.config.MonzoAccountID,
		"url":        url,
	})
	if err != nil {
		return err
	}
	if status != 200 {
		return fmt.Errorf("API did not return 200: %v %v", status, string(body))
	}

	return nil
}

// Webhook represents a Monzo webhook
type Webhook struct {
	AccountID string `json:"account_id"`
	ID        string `json:"id"`
	URL       string `json:"url"`
}

type listWebhookResponse struct {
	Webhooks []Webhook `json:"webhooks"`
}

// ListWebhooks returns a slice of webhooks currently registered.
func (g Gateway) ListWebhooks() ([]Webhook, error) {
	status, body, err := g.client.GET(fmt.Sprintf("%s/webhooks?account_id=%s", monzoAPIUrl, g.config.MonzoAccountID))
	if err != nil {
		return []Webhook{}, err
	}
	if status != 200 {
		return []Webhook{}, fmt.Errorf("API returned error: %v %v", status, body)
	}

	var response listWebhookResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return []Webhook{}, fmt.Errorf("Could not read response: %w", err)
	}

	return response.Webhooks, nil
}

// GetTransaction fetches a Transaction from an ID.
func (g Gateway) GetTransaction(id string) (Transaction, error) {
	status, body, err := g.client.GET(fmt.Sprintf("%s/transactions/%s", monzoAPIUrl, id))
	if err != nil {
		return Transaction{}, err
	}
	if status != 200 {
		return Transaction{}, fmt.Errorf("API returned error: %v %v", status, string(body))
	}

	var response transactionResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return Transaction{}, err
	}

	return response.Transaction, nil
}
