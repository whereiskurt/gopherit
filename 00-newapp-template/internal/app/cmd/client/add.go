package client

import (
	"00-newapp-template/pkg/client"
	"00-newapp-template/pkg/ui"
	"fmt"
)

// Add new gopher, and return the added gopher
func Add(a *client.Adapter, cli ui.CLI) {
	a.Config.Client.EnableLogging()

	log := a.Config.Log

	var g client.Gopher
	if a.Config.Client.ThingID == "" {
		log.Debugf("Update a Gopher")
		g = client.Gopher{
			ID:          a.Config.Client.GopherID,
			Name:        a.Config.Client.GopherName,
			Description: a.Config.Client.GopherDescription,
			Things:      map[string]client.Thing{},
		}
		a.AddGopher(g)

		output := cli.Render("GopherTable", []client.Gopher{g})
		fmt.Println(output)

	} else {

	}

	return
}
