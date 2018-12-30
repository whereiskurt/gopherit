package main

import (
	"00-newapp-template/internal"
	"00-newapp-template/internal/pkg"
)

func main() {
	config := pkg.NewConfig()
	metrics := pkg.NewMetrics(config.Metrics)

	a := internal.NewApp(config, metrics)

	a.InvokeCLI()
	return
}
