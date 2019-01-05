package client

import (
	"00-newapp-template/pkg/adapter"
	"00-newapp-template/pkg/ui"
	"encoding/json"
	"fmt"
	"strings"
)

// List uses a configured adapter to list matching gopher/things.
// Returns all matching gophers
func List(a *adapter.Adapter, cli ui.CLI) map[string]adapter.Gopher {
	a.Config.EnableClientLogging()

	log := a.Config.Log

	gophers := a.GopherThings()

	log.Debugf("Gopher Map retrieved: %+v", gophers)

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

	a.Config.Client.DumpMetrics()

	return gophers
}
