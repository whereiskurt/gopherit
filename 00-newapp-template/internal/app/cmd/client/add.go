package client

import (
	"00-newapp-template/pkg/client"
	"00-newapp-template/pkg/ui"
	"encoding/json"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"strings"
)

// Add new gopher, and return the added gopher
func Add(a *client.Adapter, cli ui.CLI) {
	a.Config.Client.EnableLogging()

	log := a.Config.Log

	if a.Config.Client.ThingID == "" {
		var g client.Gopher
		log.Debugf("Adding a Gopher from this config: %s", spew.Sprintf("%v", a.Config))
		g = client.Gopher{
			ID:          a.Config.Client.GopherID,
			Name:        a.Config.Client.GopherName,
			Description: a.Config.Client.GopherDescription,
			Things:      map[string]client.Thing{},
		}
		a.AddGopher(g)
		outputAddGopher(a, cli, g)

	} else if a.Config.Client.ThingID == "" {
		var t client.Thing
		log.Debugf("Adding a Thing from this config: %s", spew.Sprintf("%v", a.Config))
		t = client.Thing{
			ID:          a.Config.Client.ThingID,
			Name:        a.Config.Client.ThingName,
			Description: a.Config.Client.ThingDescription,
			Gopher: client.Gopher{
				ID: a.Config.Client.GopherID,
			},
		}
		a.AddThing(t)
		outputAddThing(a, cli, t)
	}

	return

}

func outputAddGopher(a *client.Adapter, cli ui.CLI, g client.Gopher) {
	var output string
	switch strings.ToLower(a.Config.Client.OutputMode) {
	case "csv":
	case "json":
		bb, _ := json.Marshal(g)
		output = string(bb)
	default:
		output = cli.Render("GopherTable", []client.Gopher{g})

	}
	fmt.Println(output)
}
func outputAddThing(a *client.Adapter, cli ui.CLI, t client.Thing) {
	var output string
	switch strings.ToLower(a.Config.Client.OutputMode) {
	case "csv":
	case "json":
		bb, _ := json.Marshal(t)
		output = string(bb)
	default:
		output = cli.Render("ThingTable", []client.Thing{t})
	}
	fmt.Println(output)
}
