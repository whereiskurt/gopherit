package client

import (
	"00-newapp-template/pkg/acme"
	"00-newapp-template/pkg/config"
	"strings"
)

// Filter is used adapter to remove unneeded results (ie. only matching Gophers)
type Filter struct {
	Config *config.Config
}

// NewFilter loops through in[] and keeps/skips matching items based on attributes.
func NewFilter(config *config.Config) (filter *Filter) {
	filter = new(Filter)
	filter.Config = config
	return
}

type gopherFilter struct {
	gg []acme.Gopher
}
type thingFilter struct {
	tt []acme.Thing
}

func (f *Filter) gophers(in []acme.Gopher) (out []acme.Gopher) {
	gopherID := f.Config.Client.GopherID
	name := f.Config.Client.GopherName

	// IDs can be comma separated and don't guaranteed just ONE match.
	if gopherID != "" {
		gg := strings.Split(gopherID, ",")
		for i := range gg {
			filter := &gopherFilter{gg: in}
			out = append(out, filter.id(gg[i])...)
		}
		if len(out) == 0 {
			return
		}
	} else {
		out = in
	}

	if name != "" {
		filter := &gopherFilter{gg: out}
		out = filter.name(name)
		if len(out) == 0 {
			return
		}
	}

	return
}
func (f *Filter) things(in []acme.Thing) (out []acme.Thing) {
	thingID := f.Config.Client.ThingID
	name := f.Config.Client.ThingName
	gopherID := f.Config.Client.GopherID

	// IDs can be comma separated and don't guaranteed just ONE match.
	if thingID != "" {
		gg := strings.Split(thingID, ",")
		for i := range gg {
			filter := &thingFilter{tt: in}
			out = append(out, filter.id(gg[i])...)
		}
		if len(out) == 0 {
			return
		}
	} else {
		out = in
	}

	if gopherID != "" {
		gg := strings.Split(thingID, ",")
		for i := range gg {
			filter := &thingFilter{tt: out}
			out = append(out, filter.gopherID(gg[i])...)
		}
		if len(out) == 0 {
			return
		}
	}

	if name != "" {
		filter := &thingFilter{tt: out}
		out = filter.name(name)
	}

	return
}

func (gf *gopherFilter) id(id string) (out []acme.Gopher) {
	for i := range gf.gg {
		if string(gf.gg[i].ID) == id {
			out = append(out, gf.gg[i])
		}
	}
	return
}
func (gf *gopherFilter) name(name string) (out []acme.Gopher) {
	for i := range gf.gg {
		if strings.Contains(strings.ToLower(gf.gg[i].Name), strings.ToLower(name)) {
			out = append(out, gf.gg[i])
		}
	}
	return
}
func (gf *thingFilter) id(id string) (out []acme.Thing) {
	for _, t := range gf.tt {
		if string(t.ID) == id {
			out = append(out, t)
		}
	}
	return
}
func (gf *thingFilter) gopherID(id string) (out []acme.Thing) {
	for _, t := range gf.tt {
		if string(t.GopherID) == id {
			out = append(out, t)
		}
	}
	return
}
func (gf *thingFilter) name(name string) (out []acme.Thing) {
	for i := range gf.tt {
		if strings.Contains(strings.ToLower(gf.tt[i].Name), strings.ToLower(name)) {
			out = append(out, gf.tt[i])
		}
	}
	return
}
