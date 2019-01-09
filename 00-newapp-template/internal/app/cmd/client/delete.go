package client

import (
	"00-newapp-template/pkg/client"
	"00-newapp-template/pkg/ui"
	"encoding/json"
	"fmt"
	"strings"
)

// Delete uses a configured adapter to delete matching gopher.
// Returns all other gophers not deleted.
func Delete(a *client.Adapter, cli ui.CLI) (gophers map[string]client.Gopher) {
	a.Config.Client.EnableLogging()
	log := a.Config.Log

	gopherID := a.Config.Client.GopherID
	thingID := a.Config.Client.ThingID

	switch {
	case gopherID != "" && thingID != "":
		tt := strings.Split(thingID, ",")
		for i := range tt {
			things := a.DeleteThing(gopherID, tt[i])
			log.Debugf("Remaining things Gopher[%s]:%d", gopherID, len(things))
			gophers := a.GopherThings()

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

	case gopherID != "":
		gg := strings.Split(gopherID, ",")
		for i := range gg {
			gophers := a.DeleteGopher(gg[i])
			log.Debugf("Gophers remaining: %d", len(gophers))
			var output string
			switch strings.ToLower(a.Config.Client.OutputMode) {
			case "csv":
			case "json":
				bb, _ := json.Marshal(gophers)
				output = string(bb)
			default:
				output = cli.Render("GopherTable", gophers)
			}
			fmt.Println(output)
		}

	case thingID != "":
		tt := strings.Split(thingID, ",")
		for _, thingID := range tt {
			gopherID = a.FindGopherByThing(thingID)
			if gopherID != "" {
				things := a.DeleteThing(gopherID, thingID)
				log.Debugf("Remaining Things for Gopher[%s]: %d", gopherID, len(things))
				var output string
				switch strings.ToLower(a.Config.Client.OutputMode) {
				case "csv":
				case "json":
					bb, _ := json.Marshal(things)
					output = string(bb)
				default:
					output = cli.Render("ThingTable", things)
				}
				fmt.Println(output)

			}
		}
	}

	things := a.GopherThings()

	a.Config.Client.DumpMetrics()
	return things
}
