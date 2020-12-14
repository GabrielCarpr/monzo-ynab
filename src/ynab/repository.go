package ynab

import (
	"fmt"

	"monzo-ynab/domain"
	"monzo-ynab/internal/config"
)

// NewRepository returns a new YNAB repository
func NewRepository(config config.Config, gateway *Gateway) *Repository {
	return &Repository{config, gateway}
}

// Repository is a YNAB repository
type Repository struct {
	config  config.Config
	gateway *Gateway
}

// newYnabTransaction converts a domain transaction to a ynab DTO transaction
func (r Repository) newYnabTransaction(t domain.Transaction) ynabTransaction {
	yt := ynabTransaction{
		AccountID: r.config.YNABAccountID,
		PayeeName: t.Payee,
		Date:      t.Date.Format(ynabDateTimeLayout),
		Amount:    t.Amount * 10,
		Memo:      t.Description,
		Cleared:   CLEARED,
		Approved:  false,
	}
	yt.generateImportID()
	return yt
}

// Store stores a transaction in YNAB
func (r Repository) Store(t domain.Transaction) error {
	ynabTrans := r.newYnabTransaction(t)

	err := r.gateway.CreateTransaction(ynabTrans)
	if err != nil {
		return fmt.Errorf("Store: %w", err)
	}
	return nil
}