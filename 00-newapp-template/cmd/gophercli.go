package main

import (
	"00-newapp-template/internal"
	"00-newapp-template/pkg/config"
	"00-newapp-template/pkg/metrics"
)

func main() {
	c := config.NewConfig()
	m := metrics.NewMetrics()

	a := internal.NewApp(c, m)

	a.InvokeCLI()
	return
}
