package commands

import (
	"fmt"
	"log"
	"monzo-ynab/domain"
	"monzo-ynab/internal/config"
)

type iYNABRepo interface {
	Store(domain.Transaction) error
}

type iMonzoRepo interface {
	List(days int) ([]domain.Transaction, error)
}

// NewSync returns a new Sync Command.
func NewSync(config config.Config, yR iYNABRepo, mR iMonzoRepo) *Sync {
	return &Sync{config, yR, mR}
}

// Sync pulls transactions from monzo over the date period and adds them into YNAB.
type Sync struct {
	config    config.Config
	ynabRepo  iYNABRepo
	monzoRepo iMonzoRepo
}

// Execute runs the command, taking the number of days to sync.
func (c Sync) Execute(days int) error {
	transactions, err := c.monzoRepo.List(days)
	if err != nil {
		return fmt.Errorf("Could not get transactions from Monzo: %w", err)
	}
	log.Printf("Retrieved %v transactions from Monzo", len(transactions))

	for _, tx := range transactions {
		if err := c.ynabRepo.Store(tx); err != nil {
			return fmt.Errorf("Could not persist transaction to YNAB: %w", err)
		}
	}
	log.Printf("Saved %v transactions to YNAB", len(transactions))
	return nil
}
