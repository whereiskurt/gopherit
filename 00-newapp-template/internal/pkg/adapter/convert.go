package adapter

import (
	"00-newapp-template/pkg/acme"
)

// Convert does need any other objects or references
type Convert struct{}

// NewConvert returns a converter, used by the adapter
func NewConvert() (convert Convert) { return }

func (convert *Convert) gophers(rawGophers []acme.Gopher) (gg map[string]Gopher) {
	gg = make(map[string]Gopher)
	for _, g := range rawGophers {
		id := string(g.ID)
		gg[id] = Gopher{
			ID:          id,
			Name:        g.Name,
			Description: g.Description,
		}
	}
	return
}

func (convert *Convert) things(rawThings []acme.Thing) (tt map[string]Thing) {
	tt = make(map[string]Thing)
	for _, t := range rawThings {
		id := string(t.ID)
		tt[id] = Thing{
			ID:          id,
			Name:        t.Name,
			Description: t.Description,
			Gopher:      Gopher{ID: string(t.GopherID)},
		}
	}
	return
}
