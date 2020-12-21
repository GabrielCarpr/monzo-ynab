package commands

import (
	"log"
	"monzo-ynab/domain"
	"monzo-ynab/internal/config"
	"monzo-ynab/ynab"
)

// NewStoreCommand constructs a Convert command.
func NewStoreCommand(c config.Config, y *ynab.Repository) *Store {
	return &Store{c, y}
}

// Store adds a new transaction to YNAB from a Monzo transaction using an ID.
type Store struct {
	config config.Config
	ynab   *ynab.Repository
}

// Execute runs the command.
func (c Store) Execute(trans domain.Transaction) error {
	err := c.ynab.Store(trans)
	if err != nil {
		return err
	}
	log.Printf("Stored Monzo transaction %s in YNAB", trans.MonzoID)

	return nil
}
