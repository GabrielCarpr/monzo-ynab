package domain

import "time"

// A transactable is a type that can be transformd into
// a domain Transaction.
type Transactable interface {
	Transaction() Transaction
}

// Transaction is the main type of transaction that needs to be synced.
type Transaction struct {
	YNABId      string
	Amount      int
	Date        time.Time
	Payee       string
	Description string
}
