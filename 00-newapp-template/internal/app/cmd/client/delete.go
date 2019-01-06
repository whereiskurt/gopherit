package client

import (
	"00-newapp-template/pkg/client"
	"00-newapp-template/pkg/ui"
	"fmt"
	"strings"
)

// Delete uses a configured adapter to delete matching gopher.
// Returns all other gophers not deleted.
func Delete(a *client.Adapter, cli ui.CLI) (gophers map[string]client.Gopher) {
	gopherID := a.Config.Client.GopherID
	thingID := a.Config.Client.ThingID
	log := a.Config.Log

	switch {
	case gopherID != "" && thingID != "":
		tt := strings.Split(thingID, ",")
		for i := range tt {
			things := a.DeleteThing(gopherID, tt[i])
			log.Debugf("Remaining things Gopher[%s]:%d", gopherID, len(things))
			gophers := a.GopherThings()
			fmt.Println(cli.Render("GopherThingsTable", gophers))
		}

	case gopherID != "":
		gg := strings.Split(gopherID, ",")
		for i := range gg {
			gophers := a.DeleteGopher(gg[i])
			log.Debugf("Gophers remaining: %d", len(gophers))
			fmt.Println(cli.Render("GopherTable", gophers))
		}

	case thingID != "":
		tt := strings.Split(thingID, ",")
		for _, thingID := range tt {
			gopherID = a.FindGopherByThing(thingID)
			if gopherID != "" {
				things := a.DeleteThing(gopherID, thingID)
				log.Debugf("Remaining Things for Gopher[%s]: %d", gopherID, len(things))
				fmt.Println(cli.Render("ThingTable", things))
			}
		}
	}

	return a.GopherThings()
}
