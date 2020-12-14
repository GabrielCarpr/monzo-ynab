package monzo

import (
	"fmt"
	"monzo-ynab/domain"
	"monzo-ynab/internal/config"
	"time"
)

// NewTransactionRepository returns a monzo TransactionRepository
func NewTransactionRepository(c config.Config, g *Gateway) *TransactionRepository {
	return &TransactionRepository{c, g}
}

// TransactionRepository stores transactions in Monzo.
type TransactionRepository struct {
	config  config.Config
	gateway *Gateway
}

// List lists all transactions since x days.
func (r TransactionRepository) List(days int) ([]domain.Transaction, error) {
	if days > 90 || days < 1 {
		return []domain.Transaction{}, fmt.Errorf("Days must be < 91 and > 0")
	}

	since := time.Now().Add(-time.Hour * 24 * time.Duration(days)).Format(time.RFC3339)
	txs, err := r.gateway.ListTransactions(since)
	if err != nil {
		return []domain.Transaction{}, err
	}

	var transactions []domain.Transaction
	for _, tx := range txs {
		transactions = append(transactions, tx.Transaction())
	}
	return transactions, nil
}
