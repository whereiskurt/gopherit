package main

import (
	"00-newapp-template/internal"
	"00-newapp-template/pkg/config"
	"00-newapp-template/pkg/metrics"
)

func main() {
	config := config.NewConfig()
	metrics := metrics.NewMetrics()

	a := internal.NewApp(config, metrics)

	a.InvokeCLI()
	return
}
