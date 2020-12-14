package main

import (
	"log"
	"monzo-ynab/domain"
	client "monzo-ynab/internal/client"
	"monzo-ynab/internal/config"
	"monzo-ynab/ynab"
	"time"

	"github.com/sarulabs/di"
)

type App struct {
	config         config.Config
	ynabRepository *ynab.Repository
}

func (a App) Run() {
	trans := domain.Transaction{
		Amount:      100000000,
		Date:        time.Now(),
		Payee:       "Charlotte",
		Description: "Heres loads of money",
	}
	err := a.ynabRepository.Store(trans)
	if err != nil {
		log.Fatalf("Failed: %s", err)
	}
}

func buildApp(config config.Config) App {
	builder, _ := di.NewBuilder()

	builder.Add(di.Def{
		Name: "ynab-gateway",
		Build: func(ctn di.Container) (interface{}, error) {
			client := client.NewClient(config.YNABToken)
			return ynab.NewGateway(config, client), nil
		},
	})

	builder.Add(di.Def{
		Name: "app",
		Build: func(ctn di.Container) (interface{}, error) {
			repo := ctn.Get("ynab-repository").(*ynab.Repository)
			return App{config, repo}, nil
		},
	})

	builder.Add(di.Def{
		Name: "ynab-repository",
		Build: func(ctn di.Container) (interface{}, error) {
			gateway := ctn.Get("ynab-gateway").(*ynab.Gateway)
			return ynab.NewRepository(config, gateway), nil
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
