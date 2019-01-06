package db

import (
	"00-newapp-template/pkg/acme"
)

// SimpleDB is two arrays holding Gophers and Things
type SimpleDB struct {
	gg []acme.Gopher
	tt []acme.Thing
}

// NewSimpleDB provides the most basic 'mock' data possible for Gophers and Things
func NewSimpleDB() (s SimpleDB) {
	s.gg = []acme.Gopher{
		{ID: "1", Name: "Gopher1", Description: "The first Gopher (#1st)"},
		{ID: "2", Name: "Gopher2", Description: "The second Gopher (#2nd)"},
		{ID: "4", Name: "Gopher4", Description: "The fourth Gopher (#4th)"},
		{ID: "8", Name: "Gopher8", Description: "The eighth Gopher (#8th)"},
	}

	s.tt = []acme.Thing{
		{ID: "1", GopherID: "1", Name: "Head", Description: "Hat"},
		{ID: "2", GopherID: "2", Name: "Head", Description: "Hat"},
		{ID: "3", GopherID: "4", Name: "Head", Description: "Hat"},
		{ID: "4", GopherID: "8", Name: "Head", Description: "Hat"},

		{ID: "5", GopherID: "1", Name: "Feet", Description: "Shoes"},
		{ID: "6", GopherID: "2", Name: "Feet", Description: "Shoes"},
		{ID: "7", GopherID: "4", Name: "Feet", Description: "Shoes"},
		{ID: "8", GopherID: "8", Name: "Feet", Description: "Shoes"},

		{ID: "9", GopherID: "1", Name: "Waist", Description: "Belt"},
		{ID: "10", GopherID: "2", Name: "Waist", Description: "Belt"},
		{ID: "11", GopherID: "4", Name: "Waist", Description: "Belt"},
		{ID: "12", GopherID: "8", Name: "Waist", Description: "Belt"},
	}
	return
}

// GopherThings returns array of acme.Things for a given Gopher ID
func (s *SimpleDB) GopherThings(gopherID string) (things []acme.Thing) {
	for _, v := range s.tt {
		if string(v.GopherID) == gopherID {
			things = append(things, v)
		}
	}
	return
}

// Gophers returns array of acme.Gophers
func (s *SimpleDB) Gophers() []acme.Gopher {
	return s.gg
}

// DeleteGopher 'cascade deleted' from gophers and things.
func (s *SimpleDB) DeleteGopher(gopherID string) {
	var gophers []acme.Gopher
	var things []acme.Thing

	for _, g := range s.gg {
		if string(g.ID) == gopherID {
			continue
		}
		gophers = append(gophers, g)
	}
	s.gg = gophers

	for _, t := range s.tt {
		if string(t.GopherID) == gopherID {
			continue
		}
		things = append(things, t)
	}
	s.tt = things
}

// UpdateGopher replaces the matching Gopher with the one passed in.
func (s *SimpleDB) UpdateGopher(newGopher acme.Gopher) {
	var gophers []acme.Gopher
	for _, g := range s.gg {
		if string(newGopher.ID) == string(g.ID) {
			gophers = append(gophers, newGopher)
			continue
		}
		gophers = append(gophers, g)
	}
	s.gg = gophers
}


func (s *SimpleDB) UpdateThing(newThing acme.Thing) {
	var things []acme.Thing

	for _, t := range s.tt {
		if (t.GopherID == newThing.GopherID) && (t.ID == newThing.ID) {
			things = append(things,newThing)
			continue
		}
		things = append(things,t)
	}
	s.tt = things
}



// DeleteThing deletes Thing that matches ID and Gopher ID
func (s *SimpleDB) DeleteThing(gopherID string, thingID string) {
	var things []acme.Thing

	for _, t := range s.tt {
		if string(t.GopherID) == gopherID && string(t.ID) == thingID {
			continue
		}
		things = append(things, t)
	}

	s.tt = things
}
