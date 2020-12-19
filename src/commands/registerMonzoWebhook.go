package commands

import (
	"fmt"
	"monzo-ynab/internal/config"
	"monzo-ynab/monzo"
)

// NewRegisterWebhook returns a RegisterWebhook command.
func NewRegisterWebhook(c config.Config, g *monzo.WebhookRepository) *RegisterMonzoWebhook {
	return &RegisterMonzoWebhook{c, g}
}

// RegisterMonzoWebhook is a command that creates a Monzo webhook.
type RegisterMonzoWebhook struct {
	config     config.Config
	repository *monzo.WebhookRepository
}

// Execute runs the command.
func (c RegisterMonzoWebhook) Execute(path string) error {
	err := c.repository.Register(path)
	if err != nil {
		return err
	}
	fmt.Printf("Monzo webhook registered")
	return nil
}
