package client

import (
	"00-newapp-template/pkg/client"
	"00-newapp-template/pkg/ui"
	"encoding/json"
	"fmt"
	"strings"
)

// List uses a configured adapter to list matching gopher/things.
// Returns all matching gophers
func List(a *client.Adapter, cli ui.CLI) map[string]client.Gopher {
	a.Config.Client.EnableLogging()

	log := a.Config.Log

	gophers := a.GopherThings()

	log.Debugf("Gopher Map retrieved: %+v", gophers)

	outputList(a, gophers, cli)

	a.Config.Client.DumpMetrics()

	return gophers
}

func outputList(a *client.Adapter, gophers map[string]client.Gopher, cli ui.CLI) {
	var output string
	switch strings.ToLower(a.Config.Client.OutputMode) {
	case "csv":
	case "json":
		bb, _ := json.Marshal(gophers)
		output = string(bb)
	default:
		output = cli.Render("GopherThingsTable", gophers)
	}
	fmt.Println(output)
}
