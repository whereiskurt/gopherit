package adapter

import (
	"00-newapp-template/internal/pkg"
	"sync"
)

// Adapter is used to call ACME services and conver them to Gopher/Things in Go structures we like.
type Adapter struct {
	Config    *pkg.Config
	Unmarshal Unmarshal
	Filter    *Filter
	Convert   Convert
	Worker    *sync.WaitGroup
}

// NewAdapter manages calls the remote services, converts the results and manages a memory/disk cache.
func NewAdapter(config *pkg.Config) (a *Adapter) {
	a = new(Adapter)
	a.Config = config
	a.Worker = new(sync.WaitGroup)
	a.Unmarshal = NewUnmarshal(config)
	a.Filter = NewFilter(config)
	a.Convert = NewConvert()

	return
}

// Gopher will return a gopher which matches the gopherID
func (a *Adapter) Gopher(gopherID string) (gopher Gopher) { return }

// Things will return all things for a gopherID
func (a *Adapter) Things(gopherID string) (things map[string]Thing) {
	rawThings := a.Unmarshal.things(gopherID)
	filtered := a.Filter.things(rawThings)
	things = a.Convert.things(filtered)
	return
}

// Gophers returns all gophers with 'things' == nil
func (a *Adapter) Gophers() (gophers map[string]Gopher) {
	rawGophers := a.Unmarshal.gophers()
	filtered := a.Filter.gophers(rawGophers)
	gophers = a.Convert.gophers(filtered)
	return
}

// GopherThings populates each gopher with their things
func (a *Adapter) GopherThings() (gophers map[string]Gopher) {
	gophers = make(map[string]Gopher)

	var matchOnThings = false

	if a.Config.Client.ThingID != "" || a.Config.Client.ThingName != "" || a.Config.Client.ThingDescription != "" {
		matchOnThings = true
	}

	gg := a.Gophers()
	for _, g := range gg {
		things := a.Things(g.ID)

		// If there are no 'things' for this gopher and we are filtering for a thing
		// don't add this 'gopher' to the results
		if len(things) == 0 && matchOnThings {
			continue
		}
		gophers[g.ID] = Gopher{
			ID:          g.ID,
			Name:        g.Name,
			Description: g.Description,
			Things:      things,
		}
	}
	return
}

// DeleteGopher will delete the matching gopherID
func (a *Adapter) DeleteGopher(gopherID string) (gophers map[string]Gopher) {
	rawGophers := a.Unmarshal.deleteGopher(gopherID)
	gophers = a.Convert.gophers(rawGophers)
	return
}

// DeleteThing will delete the Thing matching gopherID and thingID - could use FindGopherByThing instead of taking thingID
func (a *Adapter) DeleteThing(gopherID string, thingID string) (things map[string]Thing) {
	rawThings := a.Unmarshal.deleteThing(gopherID, thingID)
	things = a.Convert.things(rawThings)
	return
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
