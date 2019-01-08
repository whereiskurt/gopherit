package server_test

import (
	"00-newapp-template/pkg/acme"
	"00-newapp-template/pkg/config"
	"00-newapp-template/pkg/server"
	"bytes"
	"encoding/json"
	"github.com/davecgh/go-spew/spew"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestAddGopherHandler(t *testing.T) {
	c := config.NewConfig()

	c.Server.CacheResponse = false

	// Clean up from last run
	_ = os.RemoveAll(c.Server.CacheFolder)
	_ = os.RemoveAll(c.Client.CacheFolder)

	s := server.NewServer(c, m)

	gc := len(s.DB.Gophers())

	g := acme.Gopher{
		ID:          "42",
		Name:        "Douglas",
		Description: "Adams",
	}

	body, err := json.Marshal(g)
	if err != nil {
		t.Fatalf("could not marshal gopher: %v", err)
	}

	// justforfunc #16: unit testing HTTP servers:
	// https://www.youtube.com/watch?v=hVFEV-ieeew
	req, reqErr := http.NewRequest(http.MethodGet, "", bytes.NewBuffer(body))
	if reqErr != nil {
		t.Fatalf("could not create request: %v", err)
	}
	rec := httptest.NewRecorder()

	s.AddGopher(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	gg := s.DB.Gophers()
	if len(gg) != gc+1 || res.StatusCode != 200 {
		t.Fatalf("failed to add gopher")
	}

	t.Log(spew.Sprintf("%+v", gg))
}
func TestAddThingHandler(t *testing.T) {

	c := config.NewConfig()
	c.Server.CacheResponse = false

	// Clean up from last run
	_ = os.RemoveAll(c.Server.CacheFolder)
	_ = os.RemoveAll(c.Client.CacheFolder)

	s := server.NewServer(c, m)

	gc := len(s.DB.Things())

	g := acme.Thing{
		GopherID:    "1",
		ID:          "102",
		Name:        "Hitchers",
		Description: "Real Nice",
	}

	body, err := json.Marshal(g)
	if err != nil {
		t.Fatalf("could not marshal gopher: %v", err)
	}

	// justforfunc #16: unit testing HTTP servers:
	// https://www.youtube.com/watch?v=hVFEV-ieeew
	req, reqErr := http.NewRequest(http.MethodGet, "", bytes.NewBuffer(body))
	if reqErr != nil {
		t.Fatalf("could not create request: %v", err)
	}
	rec := httptest.NewRecorder()

	s.AddThing(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	gg := s.DB.Things()
	if len(gg) != gc+1 || res.StatusCode != 200 {
		t.Fatalf("failed to add thing")
	}

	t.Log(spew.Sprintf("%+v", gg))
}
