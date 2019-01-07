package db_test

import (
	"00-newapp-template/pkg/acme"
	"00-newapp-template/pkg/server/db"
	"encoding/json"
	"github.com/davecgh/go-spew/spew"
	"testing"
)

func TestGopherAddDelete(t *testing.T) {
	t.Parallel()

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
