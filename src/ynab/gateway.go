package ynab

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"monzo-ynab/domain"
	"monzo-ynab/internal/config"
	"net/http"
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

// ynabTransaction is a model of the YNAB API transaction object.
type ynabTransaction struct {
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
func (t *ynabTransaction) assignAccountID(id string) {
	t.AccountID = id
}

// generateImportID creates an import ID for the transaction
func (t *ynabTransaction) generateImportID() {
	formatStr := "YNAB:%v:%s:1"
	date, _ := time.Parse(ynabDateLayout, t.Date)
	datetime := date.Format(ynabDateTimeLayout)
	t.ImportID = fmt.Sprintf(formatStr, t.Amount, datetime)
}

// Transaction implements the transactable interface.
func (t ynabTransaction) Transaction() domain.Transaction {
	date, err := time.Parse(ynabDateLayout, t.Date)
	if err != nil {
		panic(err)
	}

	return domain.Transaction{
		YNABId:      t.ID,
		Amount:      t.Amount / 10,
		Date:        date,
		Payee:       t.PayeeName,
		Description: t.Memo,
	}
}

// NewGateway returns a configured, useable Gateway.
func NewGateway(config config.Config) *Gateway {
	client := http.Client{}
	return &Gateway{config, client}
}

// Gateway is the Gateway over the YNAB API
type Gateway struct {
	config config.Config
	client http.Client
}

// CreateTransaction posts a transaction to the YNAB API
func (g Gateway) CreateTransaction(transaction ynabTransaction) error {
	goBody := map[string]ynabTransaction{"transaction": transaction}
	body, err := json.Marshal(goBody)
	if err != nil {
		return fmt.Errorf("CreateTransaction: %w", err)
	}

	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s/budgets/%s/transactions", ynabAPI, g.config.YNABBudgetID), // URL
		bytes.NewBuffer(body), // Body
	)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", g.config.YNABToken))
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		return fmt.Errorf("CreateTransaction: %w")
	}

	resp, err := g.client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return fmt.Errorf("CreateTransaction: %w", err)
	}
	if resp.StatusCode == 201 {
		log.Printf("Added transaction %s", transaction.Memo)
		return nil
	}
	if resp.StatusCode == 400 {
		return fmt.Errorf("CreateTransaction: Bad request")
	}
	if resp.StatusCode == 409 {
		log.Printf("Transaction already exists")
		return nil // The transaction already exists.
	}
	bod, err := ioutil.ReadAll(resp.Body)
	log.Print(string(bod))
	return fmt.Errorf("CreateTransaction: Unknown response %v", resp.StatusCode)
}
