package adapter

import (
	"00-newapp-template/internal/pkg"
	"00-newapp-template/pkg/cache"
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"sync"
)

// CacheLabel is the type for where to store the response
type CacheLabel string
func (c CacheLabel) String() string {
	return "adapter/" + string(c)
}

// Adapter is used to call ACME services and convert them to Gopher/Things in Go structures we like.
type Adapter struct {
	Config    *pkg.Config
	Unmarshal Unmarshal
	Filter    *Filter
	Convert   Convert
	Worker    *sync.WaitGroup
	DiskCache *cache.Disk
}

// NewAdapter manages calls the remote services, converts the results and manages a memory/disk cache.
func NewAdapter(config *pkg.Config) (a *Adapter) {
	a = new(Adapter)
	a.Config = config
	a.Worker = new(sync.WaitGroup)
	a.Unmarshal = NewUnmarshal(config)
	a.Filter = NewFilter(config)
	a.Convert = NewConvert()
	if a.Config.Client.CacheResponse {
		a.DiskCache = cache.NewDisk(a.Config.Client.CacheFolder, a.Config.Client.CacheKey, a.Config.Client.CacheKey != "")
	}

	return
}

// Things will return all things for a gopherID
func (a *Adapter) Things(gopherID string) map[string]Thing {
	rawThings := a.Unmarshal.things(gopherID)
	filtered := a.Filter.things(rawThings)
	things := a.Convert.things(filtered)

	label := CacheLabel(fmt.Sprintf("Things/Gopher.%s", gopherID))
	a.CacheStore(label, &things )

	return things
}

// Gophers returns all gophers with 'things' == nil
func (a *Adapter) Gophers() map[string]Gopher {
	rawGophers := a.Unmarshal.gophers()
	filtered := a.Filter.gophers(rawGophers)
	gophers := a.Convert.gophers(filtered)

	a.CacheStore(CacheLabel("Gophers"), &gophers )

	return gophers
}

// GopherThings populates each gopher with their things
func (a *Adapter) GopherThings() map[string]Gopher {
	var matchOnThings = false
	if a.Config.Client.ThingID != "" || a.Config.Client.ThingName != "" || a.Config.Client.ThingDescription != "" {
		matchOnThings = true
	}

	gopherThings := make(map[string]Gopher)

	gg := a.Gophers()
	for _, g := range gg {
		things := a.Things(g.ID)

		// If there are no 'things' for this gopher and we are filtering for a thing
		// don't add this 'gopher' to the results
		if len(things) == 0 && matchOnThings {
			continue
		}
		gopherThings[g.ID] = Gopher{
			ID:          g.ID,
			Name:        g.Name,
			Description: g.Description,
			Things:      things,
		}
	}

	a.CacheStore(CacheLabel("GopherThings"), &gopherThings)
	return gopherThings
}

// DeleteGopher will delete the matching gopherID
func (a *Adapter) DeleteGopher(gopherID string) map[string]Gopher {
	rawGophers := a.Unmarshal.deleteGopher(gopherID)
	gophers := a.Convert.gophers(rawGophers)
	return gophers
}

// DeleteThing will delete the Thing matching gopherID and thingID - could use FindGopherByThing instead of taking thingID
func (a *Adapter) DeleteThing(gopherID string, thingID string) map[string]Thing {
	rawThings := a.Unmarshal.deleteThing(gopherID, thingID)
	things := a.Convert.things(rawThings)
	return things
}

// FindGopherByThing returns the Gopher ID to the associated Thing by ID.
func (a *Adapter) FindGopherByThing(thingID string) (gopherID string) {
	gophers := a.GopherThings()
	for g := range gophers {
		for t := range gophers[g].Things {
			if string(gophers[g].Things[t].ID) == thingID {
				gopherID = gophers[g].ID
				return
			}
		}
	}
	return
}

// UpdateGopher is not implemented yet!
func (a *Adapter) UpdateGopher(newGopher Gopher) (gopher Gopher) { return }

// UpdateThing is not implemented yet!
func (a *Adapter) UpdateThing(newThing Thing) (thing Thing) { return }

func (a *Adapter) CacheStore(name CacheLabel, obj interface{}) {
	json, err := json.Marshal(obj)
	if err == nil {
		a.DiskCache.Store(fmt.Sprintf("%s.json", name), Prettify(json))
	}
}

func Prettify(json []byte) []byte {
	jq, err := exec.LookPath("jq")
	if err == nil {
		var pretty bytes.Buffer
		cmd := exec.Command(jq, ".")
		cmd.Stdin = strings.NewReader(string(json))
		cmd.Stdout = &pretty
		err := cmd.Run()
		if err == nil {
			json = []byte(pretty.String())
		}
	}
	return json
}