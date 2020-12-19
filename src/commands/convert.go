package commands

import (
	"log"
	"monzo-ynab/internal/config"
	"monzo-ynab/monzo"
	"monzo-ynab/ynab"
)

// NewConvertCommand constructs a Convert command.
func NewConvertCommand(c config.Config, y *ynab.Repository, m *monzo.TransactionRepository) *Convert {
	return &Convert{c, y, m}
}

// Convert adds a new transaction to YNAB from a Monzo transaction using an ID.
type Convert struct {
	config config.Config
	ynab   *ynab.Repository
	monzo  *monzo.TransactionRepository
}

// Execute runs the command.
func (c Convert) Execute(id string) error {
	tx, err := c.monzo.Get(id)
	if err != nil {
		return err
	}
	log.Printf("Got transaction %s from Monzo", id)

	err = c.ynab.Store(tx)
	if err != nil {
		return err
	}
	log.Printf("Stored Monzo transaction %s in YNAB", id)

	return nil
}
