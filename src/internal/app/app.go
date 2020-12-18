package app

import (
	"log"
	"monzo-ynab/commands"
	client "monzo-ynab/internal/client"
	"monzo-ynab/internal/config"
	"monzo-ynab/monzo"
	"monzo-ynab/ynab"

	di "github.com/sarulabs/di/v2"
)

// App is the main App object.
type App struct {
	config   config.Config
	commands *commands.Commands
}

// Run starts the app.
func (a App) Run() {
	err := a.commands.RegisterWebhook.Execute("/test")
	if err != nil {
		log.Fatal(err)
	}
}

// BuildApp returns a DI container.
func BuildApp(config config.Config) *di.Builder {
	builder, _ := di.NewBuilder()

	builder.Add(di.Def{
		Name: "app",
		Build: func(ctn di.Container) (interface{}, error) {
			repo := ctn.Get("commands").(*commands.Commands)
			return App{config, repo}, nil
		},
	})

	builder.Add(di.Def{
		Name: "ynab-gateway",
		Build: func(ctn di.Container) (interface{}, error) {
			client := ctn.Get("ynab-client").(client.IClient)
			return ynab.NewGateway(config, client), nil
		},
	})

	builder.Add(di.Def{
		Name: "ynab-repository",
		Build: func(ctn di.Container) (interface{}, error) {
			gateway := ctn.Get("ynab-gateway").(*ynab.Gateway)
			return ynab.NewRepository(config, gateway), nil
		},
	})

	builder.Add(di.Def{
		Name: "monzo-client",
		Build: func(ctn di.Container) (interface{}, error) {
			return client.NewClient(config.MonzoAccessToken), nil
		},
	})

	builder.Add(di.Def{
		Name: "monzo-gateway",
		Build: func(ctn di.Container) (interface{}, error) {
			client := ctn.Get("monzo-client").(client.IClient)
			return monzo.NewGateway(config, client), nil
		},
	})

	builder.Add(di.Def{
		Name: "monzo-repository",
		Build: func(ctn di.Container) (interface{}, error) {
			gateway := ctn.Get("monzo-gateway").(*monzo.Gateway)
			return monzo.NewTransactionRepository(config, gateway), nil
		},
	})

	builder.Add(di.Def{
		Name: "sync-command",
		Build: func(ctn di.Container) (interface{}, error) {
			ynabRepo := ctn.Get("ynab-repository").(*ynab.Repository)
			monzoRepo := ctn.Get("monzo-repository").(*monzo.TransactionRepository)
			return commands.NewSync(config, ynabRepo, monzoRepo), nil
		},
	})

	builder.Add(di.Def{
		Name: "register-webhook-command",
		Build: func(ctn di.Container) (interface{}, error) {
			monzoGateway := ctn.Get("monzo-webhook-repository").(*monzo.WebhookRepository)
			return commands.NewRegisterWebhook(config, monzoGateway), nil
		},
	})

	builder.Add(di.Def{
		Name: "monzo-webhook-repository",
		Build: func(ctn di.Container) (interface{}, error) {
			monzoGateway := ctn.Get("monzo-gateway").(*monzo.Gateway)
			return monzo.NewWebhookRepository(config, monzoGateway), nil
		},
	})

	builder.Add(di.Def{
		Name: "commands",
		Build: func(ctn di.Container) (interface{}, error) {
			return &commands.Commands{
				Sync:            ctn.Get("sync-command").(*commands.Sync),
				RegisterWebhook: ctn.Get("register-webhook-command").(*commands.RegisterWebhook),
			}, nil
		},
	})

	builder.Add(di.Def{
		Name: "ynab-client",
		Build: func(ctn di.Container) (interface{}, error) {
			return client.NewClient(config.YNABToken), nil
		},
	})

	return builder
}
