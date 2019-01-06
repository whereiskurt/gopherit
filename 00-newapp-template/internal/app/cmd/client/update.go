package client

import (
	"00-newapp-template/pkg/client"
	"00-newapp-template/pkg/ui"
)

// Update takes command line parameters to update a gopher/thing depending
func Update(a *client.Adapter, cli ui.CLI) map[string]client.Gopher {
	a.Config.Client.EnableLogging()

	log := a.Config.Log
	log.Debugf("Client Update")

	if a.Config.Client.ThingID == "" {
		var g = client.Gopher{
			ID:          a.Config.Client.GopherID,
			Name:        a.Config.Client.GopherName,
			Description: a.Config.Client.GopherDescription,
			Things:      map[string]client.Thing{},
		}
		a.UpdateGopher(g)
	} else {

		var t = client.Thing{
			ID:          a.Config.Client.ThingID,
			Name:        a.Config.Client.ThingName,
			Description: a.Config.Client.ThingDescription,
		}
		a.UpdateThing(t)
	}

	return a.GopherThings()
}
