package cli

import (
	"log"
	"monzo-ynab/commands"

	"monzo-ynab/internal/config"

	"github.com/AlecAivazis/survey/v2"
)

type answers struct {
	YNABToken     string `survey:"ynab-token"`
	YNABAccountID string `survey:"ynab-account"`
	YNABBudgetID  string `survey:"ynab-budget"`

	MonzoAccountID   string `survey:"monzo-account"`
	MonzoAccessToken string `survey:"monzo-token"`

	BaseURL string `survey:"base-url"`

	SetupMonzo  bool `survey:"setup-monzo"`
	InitialSync bool `survey:"initial-sync"`
	SyncDays    int  `survey:"sync-days"`
}

func (a answers) config() config.Config {
	return config.Config{
		YNABToken:        a.YNABToken,
		YNABAccountID:    a.YNABAccountID,
		YNABBudgetID:     a.YNABBudgetID,
		MonzoAccountID:   a.MonzoAccountID,
		MonzoAccessToken: a.MonzoAccessToken,
		BaseURL:          a.BaseURL,
	}
}

// NewInstaller returns a fresh installer
func NewInstaller(c *commands.Commands) *Installer {
	return &Installer{commands: c}
}

// Installer sets the app up on the machine.
type Installer struct {
	answers answers

	commands *commands.Commands
}

// Install runs an interactive CLI to configure the app.
func (i *Installer) Install() {
	firstQuestions := []*survey.Question{
		{
			Name:   "ynab-token",
			Prompt: &survey.Input{Message: "What's your YNAB personal access token?"},
		},
		{
			Name:   "ynab-account",
			Prompt: &survey.Input{Message: "What's your YNAB account ID?"},
		},
		{
			Name:   "ynab-budget",
			Prompt: &survey.Input{Message: "What's your YNAB budget ID?"},
		},

		{
			Name:   "setup-monzo",
			Prompt: &survey.Confirm{Message: "Do you want to connect Monzo"},
		},
	}

	secondQuestions := []*survey.Question{
		{
			Name:   "base-url",
			Prompt: &survey.Input{Message: "What URL should Monzo webhooks be sent to?"},
		},
		{
			Name:   "monzo-account",
			Prompt: &survey.Input{Message: "What's your Monzo account ID?"},
		},
		{
			Name: "monzo-token",
			Prompt: &survey.Input{
				Message: "What's your Monzo access token?",
				Help:    "This will only be used to setup webhooks and sync",
			},
		},
	}

	err := survey.Ask(firstQuestions, &i.answers)
	if err != nil {
		log.Fatal(err)
	}

	if i.answers.SetupMonzo {
		err = survey.Ask(secondQuestions, &i.answers)
		if err != nil {
			log.Fatal(err)
		}
	}

	config := i.answers.config()
	config.Persist()

	log.Print("Installation complete")
}
