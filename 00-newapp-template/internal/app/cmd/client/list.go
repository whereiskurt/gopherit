package client

import (
	"00-newapp-template/internal/pkg/adapter"
	"00-newapp-template/internal/pkg/ui"
	"encoding/json"
	"fmt"
	"strings"
)

// List uses a configured adapter to list matching gopher/things.
// Returns all matching gophers
func List(adapter *adapter.Adapter, cli ui.CLI) (gophers map[string]adapter.Gopher) {
	log := adapter.Config.Log

	gophers = adapter.GopherThings()

	log.Debugf("Gropher Map retrieved: %+v", gophers)

	var output string
	switch strings.ToLower(adapter.Config.Client.OutputMode) {
	case "csv":
	case "json":
		bb, _ := json.Marshal(gophers)
		output = string(bb)
	default:
		output = cli.Render("GopherThingsTable", gophers)
	}

	fmt.Println(output)

	return
}
