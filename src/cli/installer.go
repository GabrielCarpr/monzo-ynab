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

func NewInstaller(c *commands.Commands) *Installer {
	return &Installer{commands: c}
}

type Installer struct {
	answers answers

	commands *commands.Commands
}

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
			Name:   "monzo-account",
			Prompt: &survey.Input{Message: "What's your Monzo account ID?"},
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
			Name: "monzo-token",
			Prompt: &survey.Input{
				Message: "What's your Monzo access token?",
				Help:    "This will only be used to setup webhooks and sync",
			},
		},
		{
			Name:   "initial-sync",
			Prompt: &survey.Confirm{Message: "Do you want perform a sync now?"},
		},
	}

	thirdQuestions := []*survey.Question{
		{
			Name:   "sync-days",
			Prompt: &survey.Input{Message: "How many previous days do you want to sync?"},
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

	if i.answers.InitialSync {
		err := survey.Ask(thirdQuestions, &i.answers)
		if err != nil {
			log.Fatal(err)
		}
	}

	config := i.answers.config()
	config.Persist()
	log.Print("Configuration complete")

	if i.answers.SetupMonzo {
		log.Print("Registering Monzo webhook...")
		err = i.commands.RegisterMonzoWebhook.Execute("/events/monzo/")
		if err != nil {
			log.Fatal(err)
		}
	}

	if i.answers.InitialSync {
		log.Print("Starting sync...")
		err = i.commands.Sync.Execute(i.answers.SyncDays)
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Print("Installation complete")
}
