package server_test

import (
	"00-newapp-template/pkg/acme"
	"00-newapp-template/pkg/config"
	"00-newapp-template/pkg/server"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestGophersGet(t *testing.T) {
	c := config.NewConfig()

	_ = os.RemoveAll(c.Server.CacheFolder)
	_ = os.RemoveAll(c.Client.CacheFolder)

	s := server.NewServer(c, m)
	s.EnableDefaultRouter()

	srv := httptest.NewServer(s.Router)
	defer srv.Close()

	url := fmt.Sprintf("%s/gophers", srv.URL)
	res, err := http.Get(url)
	if err != nil {
		t.Fatalf("could not send GET /gophers request: %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("GET /gophers returned: %d", res.StatusCode)
	}
	defer res.Body.Close()
	bb, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("couldn't read response body: %v", err)
	}

	var gg []acme.Gopher
	err = json.Unmarshal(bb, &gg)
	t.Log(spew.Sprintf("%+v", gg))

	return
}
func TestGopherPut(t *testing.T) {
	c := config.NewConfig()

	_ = os.RemoveAll(c.Server.CacheFolder)
	_ = os.RemoveAll(c.Client.CacheFolder)

	s := server.NewServer(c, m)
	s.EnableDefaultRouter()

	srv := httptest.NewServer(s.Router)
	defer srv.Close()

	url := fmt.Sprintf("%s/gophers", srv.URL)

	g := acme.Gopher{
		ID:          "42",
		Name:        "Douglas",
		Description: "Adams",
	}

	body, err := json.Marshal(g)
	if err != nil {
		t.Fatalf("could not marshal gopher for test: %v", err)
	}

	req, _ := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(body))
	client := http.Client{Timeout: time.Second * 5}
	res, err := client.Do(req)
	if err != nil {
		t.Fatalf("couldn't PUT gopher: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("PUT gopher failed: %d", res.StatusCode)
	}
	bb, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("couldn't read response body: %v", err)
	}
	err = json.Unmarshal(bb, g)
	t.Log(spew.Sprintf("%+v", g))

	return
}
