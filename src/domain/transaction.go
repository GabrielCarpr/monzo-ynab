package domain

import "time"

// Transactable is a type that can be transformed into
// a domain Transaction.
type Transactable interface {
	Transaction() Transaction
}

// Transaction is the main type of transaction that needs to be synced.
type Transaction struct {
	MonzoID     string
	Amount      int
	Date        time.Time
	Payee       string
	Description string
}
