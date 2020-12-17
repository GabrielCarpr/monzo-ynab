package main

import (
	"monzo-ynab/internal/app"
	"monzo-ynab/internal/config"
)

func main() {
	config := config.NewConfig()
	builder := app.BuildApp(config)
	ctn := builder.Build()
	app := ctn.Get("app").(app.App)
	app.Run()
}
