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

// monzoTransaction represents a Monzo API DTO (or, just the fields we need.)
type monzoTransaction struct {
	Description string `json:"description"`
	Amount      int    `json:"amount"`
	Settled     string `json:"settled"`
}

type transactionResponse struct {
	Transactions []monzoTransaction `json:"transactions"`
}

func (t monzoTransaction) Transaction() domain.Transaction {
	date, _ := time.Parse(time.RFC3339, t.Settled)
	return domain.Transaction{
		Amount:      t.Amount,
		Date:        date,
		Payee:       t.Description,
		Description: t.Description,
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
func (g Gateway) ListTransactions(since string) ([]monzoTransaction, error) {
	status, body, err := g.client.GET(fmt.Sprintf("%s/transactions?account_id=%s&since=%s", monzoAPIUrl, g.config.MonzoAccountID, since))
	if err != nil {
		return []monzoTransaction{}, err
	}

	if status != 200 {
		return []monzoTransaction{}, fmt.Errorf("API did not return 200: %v %v", status, string(body))
	}

	var tr transactionResponse
	if err := json.Unmarshal(body, &tr); err != nil {
		return tr.Transactions, fmt.Errorf("Failed reading Monzo API response: %w", err)
	}

	return tr.Transactions, nil
}
