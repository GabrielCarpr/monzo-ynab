package main

import (
	"monzo-ynab/commands"
	client "monzo-ynab/internal/client"
	"monzo-ynab/internal/config"
	"monzo-ynab/monzo"
	"monzo-ynab/ynab"

	"github.com/sarulabs/di"
)

type App struct {
	config   config.Config
	commands *commands.Commands
}

func (a App) Run() {
	a.commands.Sync.Execute(10)
}

func buildApp(config config.Config) App {
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
			client := client.NewClient(config.YNABToken)
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
		Name: "monzo-gateway",
		Build: func(ctn di.Container) (interface{}, error) {
			client := client.NewClient(config.MonzoAccessToken)
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
		Name: "commands",
		Build: func(ctn di.Container) (interface{}, error) {
			return &commands.Commands{
				Sync: ctn.Get("sync-command").(*commands.Sync),
			}, nil
		},
	})

	ctn := builder.Build()
	app := ctn.Get("app").(App)
	return app
}

func main() {
	config := config.NewConfig()
	app := buildApp(config)
	app.Run()
}
