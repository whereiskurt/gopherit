package db_test

import (
	"00-newapp-template/pkg/acme"
	"00-newapp-template/pkg/server/db"
	"encoding/json"
	"github.com/davecgh/go-spew/spew"
	"testing"
)

func TestGopherAddDelete(t *testing.T) {

	d := db.NewSimpleDB()

	c := len(d.Gophers())
	d.AddGopher(acme.Gopher{
		ID:          json.Number("1960"),
		Description: "FromAdd",
		Name:        "Gophigur",
	})

	t.Log(spew.Sprintf("%+v", d.Gophers()))
	if c == len(d.Gophers()) {
		t.Fail()
	}

	d.DeleteGopher("1960")
	if c != len(d.Gophers()) {
		t.Fail()
	}
	t.Log(spew.Sprintf("%+v", d.Gophers()))

}
func TestThingAddDelete(t *testing.T) {

	d := db.NewSimpleDB()

	c := len(d.Things())
	d.AddThing(acme.Thing{
		ID:          json.Number("1979"),
		GopherID:    "1",
		Description: "NewThingerDesc",
		Name:        "NewThingsName",
	})

	t.Log(spew.Sprintf("%+v", d.Things()))
	if c == len(d.Things()) {
		t.Fail()
	}

	d.DeleteThing("1", "1979")
	if c != len(d.Things()) {
		t.Fail()
	}
	t.Log(spew.Sprintf("%+v", d.Things()))

}
