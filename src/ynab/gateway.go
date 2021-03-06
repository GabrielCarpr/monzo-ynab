package ynab

import (
	"fmt"
	"log"
	"monzo-ynab/domain"
	client "monzo-ynab/internal/client"
	"monzo-ynab/internal/config"
	"time"
)

type cleared string

// Enum for cleared value
const (
	CLEARED    cleared = "cleared"
	UNCLEARED  cleared = "uncleared"
	RECONCILED cleared = "reconciled"
)

const ynabAPI = "https://api.youneedabudget.com/v1"
const ynabDateLayout = "2006-01-02"
const ynabDateTimeLayout = ynabDateLayout + "T15:04"

// Transaction is a model of the YNAB API transaction object.
type Transaction struct {
	ID         string  `json:"id,omitempty"`
	AccountID  string  `json:"account_id"`
	PayeeID    *string `json:"payee_id"`
	PayeeName  string  `json:"payee_name,omitempty"`
	CategoryID *string `json:"category_id"`
	Date       string  `json:"date"`
	Amount     int     `json:"amount"`
	Memo       string  `json:"memo"`
	Cleared    cleared `json:"cleared"`
	Approved   bool    `json:"approved"`
	ImportID   string  `json:"import_id,omitempty"`
}

// AssignAccountID sets the account ID to sync to
func (t *Transaction) assignAccountID(id string) {
	t.AccountID = id
}

// generateImportID creates an import ID for the transaction
func (t *Transaction) generateImportID() {
	formatStr := "YNAB:%v:%s:1"
	date, _ := time.Parse(ynabDateTimeLayout, t.Date)
	datetime := date.Format(ynabDateLayout)
	t.ImportID = fmt.Sprintf(formatStr, t.Amount, datetime)
}

// Transaction implements the transactable interface.
func (t Transaction) Transaction() domain.Transaction {
	date, err := time.Parse(ynabDateLayout, t.Date)
	if err != nil {
		panic(err)
	}

	return domain.Transaction{
		Amount:      t.Amount / 10,
		Date:        date,
		Payee:       t.PayeeName,
		Description: t.Memo,
	}
}

// NewGateway returns a configured, useable Gateway.
func NewGateway(config config.Config, c client.IClient) *Gateway {
	return &Gateway{config, c}
}

// Gateway is the Gateway over the YNAB API
type Gateway struct {
	config config.Config
	client client.IClient
}

// CreateTransaction posts a transaction to the YNAB API
func (g Gateway) CreateTransaction(transaction Transaction) error {
	goBody := client.JSONBody{"transaction": transaction}

	status, body, err := g.client.POST(
		fmt.Sprintf("%s/budgets/%s/transactions", ynabAPI, g.config.YNABBudgetID),
		goBody,
	)
	if err != nil {
		return err
	}

	if status == 201 {
		log.Printf("Added transaction %s (%s)", transaction.Memo, transaction.ImportID)
		return nil
	}
	if status == 400 {
		return fmt.Errorf("CreateTransaction: Bad request - %v", string(body))
	}
	if status == 409 {
		log.Printf("Transaction already exists")
		return nil // The transaction already exists.
	}
	return fmt.Errorf("CreateTransaction: Unknown response %v", status)
}
