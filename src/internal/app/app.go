package app

import (
	"log"
	"monzo-ynab/commands"
	client "monzo-ynab/internal/client"
	"monzo-ynab/internal/config"
	"monzo-ynab/monzo"
	"monzo-ynab/rest"
	"monzo-ynab/ynab"
	"net/http"

	di "github.com/sarulabs/di/v2"
)

// App is the main App object.
type App struct {
	config   config.Config
	commands *commands.Commands
	rest     *rest.Handler
}

// Run starts the app.
func (a App) Run() {
	err := a.commands.RegisterMonzoWebhook.Execute("/events/monzo")
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(http.ListenAndServe(":8080", a.rest))
}

// BuildApp returns a DI container.
func BuildApp(config config.Config) *di.Builder {
	builder, _ := di.NewBuilder()

	builder.Add(di.Def{
		Name: "app",
		Build: func(ctn di.Container) (interface{}, error) {
			repo := ctn.Get("commands").(*commands.Commands)
			rest := ctn.Get("rest-handler").(*rest.Handler)
			return App{config, repo, rest}, nil
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
		Name: "store-transaction-command",
		Build: func(ctn di.Container) (interface{}, error) {
			ynabRepo := ctn.Get("ynab-repository").(*ynab.Repository)
			return commands.NewStoreCommand(config, ynabRepo), nil
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
				Sync:                 ctn.Get("sync-command").(*commands.Sync),
				RegisterMonzoWebhook: ctn.Get("register-webhook-command").(*commands.RegisterMonzoWebhook),
				Store:                ctn.Get("store-transaction-command").(*commands.Store),
			}, nil
		},
	})

	builder.Add(di.Def{
		Name: "ynab-client",
		Build: func(ctn di.Container) (interface{}, error) {
			return client.NewClient(config.YNABToken), nil
		},
	})

	builder.Add(di.Def{
		Name: "rest-handler",
		Build: func(ctn di.Container) (interface{}, error) {
			commands := ctn.Get("commands").(*commands.Commands)
			return rest.NewHandler(config, commands), nil
		},
	})

	return builder
}
