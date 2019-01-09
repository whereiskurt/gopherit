package client_test

import (
	"00-newapp-template/internal/app/cmd/client"
	"00-newapp-template/internal/app/cmd/server"
	pkgclient "00-newapp-template/pkg/client"
	"00-newapp-template/pkg/config"
	"00-newapp-template/pkg/metrics"
	"00-newapp-template/pkg/ui"
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

var m = metrics.NewMetrics()

func TestClientList(t *testing.T) {

	c := config.NewConfig()
	SetupConfig(c)

	adapter := pkgclient.NewAdapter(c, m)
	cli := ui.NewCLI(c)

	go server.Start(c, m)

	jsonOutput := captureList(adapter, cli)

	if !strings.HasPrefix(jsonOutput, `{"1":{"id":"1","name":"Gopher1","description":"The first Gopher (#1st)","things":{"1":{"gopher":{"id":"1","name":"","description":"","things":null},"id":"1","name":"Head","description":"Hat"},"5":{"gopher":{"id":"1","name":"","description":"","things":null},"id":"5","name":"Feet","description":"Shoes"},"9":{"gopher":{"id":"1","name":"","description":"","things":null},"id":"9","name":"Waist","description":"Belt"}}},"2":{"id":"2","name":"Gopher2","description":"The second Gopher (#2nd)","things":{"10":{"gopher":{"id":"2","name":"","description":"","things":null},"id":"10","name":"Waist","description":"Belt"},"2":{"gopher":{"id":"2","name":"","description":"","things":null},"id":"2","name":"Head","description":"Hat"},"6":{"gopher":{"id":"2","name":"","description":"","things":null},"id":"6","name":"Feet","description":"Shoes"}}},"4":{"id":"4","name":"Gopher4","description":"The fourth Gopher (#4th)","things":{"11":{"gopher":{"id":"4","name":"","description":"","things":null},"id":"11","name":"Waist","description":"Belt"},"3":{"gopher":{"id":"4","name":"","description":"","things":null},"id":"3","name":"Head","description":"Hat"},"7":{"gopher":{"id":"4","name":"","description":"","things":null},"id":"7","name":"Feet","description":"Shoes"}}},"8":{"id":"8","name":"Gopher8","description":"The eighth Gopher (#8th)","things":{"12":{"gopher":{"id":"8","name":"","description":"","things":null},"id":"12","name":"Waist","description":"Belt"},"4":{"gopher":{"id":"8","name":"","description":"","things":null},"id":"4","name":"Head","description":"Hat"},"8":{"gopher":{"id":"8","name":"","description":"","things":null},"id":"8","name":"Feet","description":"Shoes"}}}}`) {
		t.Fatalf("Unexpected output:%s", jsonOutput)
	}

	server.Stop(c, m)
}

func TestClientAdd(t *testing.T) {

	c := config.NewConfig()
	SetupConfig(c)

	adapter := pkgclient.NewAdapter(c, m)
	cli := ui.NewCLI(c)

	go server.Start(c, m)

	c.Client.GopherID = "42"
	c.Client.GopherName = "Dougy Addy"
	c.Client.GopherDescription = "Funny books"

	jsonOutput := captureAdd(adapter, cli)

	if !strings.HasPrefix(jsonOutput, `{"id":"42","name":"Dougy Addy","description":"Funny books","things":{}}`) {
		t.Fatalf("Unexpected output:%s", jsonOutput)
	}

	server.Stop(c, m)
}
func TestClientDeleteGopher(t *testing.T) {

	c := config.NewConfig()
	SetupConfig(c)

	adapter := pkgclient.NewAdapter(c, m)
	cli := ui.NewCLI(c)

	go server.Start(c, m)

	c.Client.GopherID = "1"
	jsonOutput := captureDelete(adapter, cli)

	if !strings.HasPrefix(jsonOutput, `{"2":{"id":"2","name":"Gopher2","description":"The second Gopher (#2nd)","things":null},"4":{"id":"4","name":"Gopher4","description":"The fourth Gopher (#4th)","things":null},"8":{"id":"8","name":"Gopher8","description":"The eighth Gopher (#8th)","things":null}}`) {
		t.Fatalf("Unexpected output:%s", jsonOutput)
	}

	server.Stop(c, m)
}
func TestClientDeleteThing(t *testing.T) {

	c := config.NewConfig()
	SetupConfig(c)

	adapter := pkgclient.NewAdapter(c, m)
	cli := ui.NewCLI(c)

	go server.Start(c, m)

	c.Client.GopherID = "1"
	c.Client.ThingID = "1,5,9"
	jsonOutput := captureDelete(adapter, cli)
	if !strings.HasPrefix(jsonOutput, `{"1":{"id":"1","name":"Gopher1","description":"The first Gopher (#1st)","things":{"5":{"gopher":{"id":"1","name":"","description":"","things":null},"id":"5","name":"Feet","description":"Shoes"},"9":{"gopher":{"id":"1","name":"","description":"","things":null},"id":"9","name":"Waist","description":"Belt"}}}}`) {
		t.Fatalf("Unexpected output:%s", jsonOutput)
	}

	c.Client.GopherID = ""
	c.Client.ThingID = ""
	jsonOutput = captureList(adapter, cli)
	if !strings.HasPrefix(jsonOutput, `{"1":{"id":"1","name":"Gopher1","description":"The first Gopher (#1st)","things":{}},"2":{"id":"2","name":"Gopher2","description":"The second Gopher (#2nd)","things":{"10":{"gopher":{"id":"2","name":"","description":"","things":null},"id":"10","name":"Waist","description":"Belt"},"2":{"gopher":{"id":"2","name":"","description":"","things":null},"id":"2","name":"Head","description":"Hat"},"6":{"gopher":{"id":"2","name":"","description":"","things":null},"id":"6","name":"Feet","description":"Shoes"}}},"4":{"id":"4","name":"Gopher4","description":"The fourth Gopher (#4th)","things":{"11":{"gopher":{"id":"4","name":"","description":"","things":null},"id":"11","name":"Waist","description":"Belt"},"3":{"gopher":{"id":"4","name":"","description":"","things":null},"id":"3","name":"Head","description":"Hat"},"7":{"gopher":{"id":"4","name":"","description":"","things":null},"id":"7","name":"Feet","description":"Shoes"}}},"8":{"id":"8","name":"Gopher8","description":"The eighth Gopher (#8th)","things":{"12":{"gopher":{"id":"8","name":"","description":"","things":null},"id":"12","name":"Waist","description":"Belt"},"4":{"gopher":{"id":"8","name":"","description":"","things":null},"id":"4","name":"Head","description":"Hat"},"8":{"gopher":{"id":"8","name":"","description":"","things":null},"id":"8","name":"Feet","description":"Shoes"}}}}`) {
		t.Fatalf("Unexpected output:%s", jsonOutput)
	}

	c.Client.ThingID = "2"
	jsonOutput = captureDelete(adapter, cli)
	if !strings.HasPrefix(jsonOutput, `{"10":{"gopher":{"id":"2","name":"","description":"","things":null},"id":"10","name":"Waist","description":"Belt"},"6":{"gopher":{"id":"2","name":"","description":"","things":null},"id":"6","name":"Feet","description":"Shoes"}}`) {
		t.Fatalf("Unexpected output:%s", jsonOutput)
	}

	c.Client.ThingID = ""
	jsonOutput = captureList(adapter, cli)
	if !strings.HasPrefix(jsonOutput, `{"1":{"id":"1","name":"Gopher1","description":"The first Gopher (#1st)","things":{}},"2":{"id":"2","name":"Gopher2","description":"The second Gopher (#2nd)","things":{"10":{"gopher":{"id":"2","name":"","description":"","things":null},"id":"10","name":"Waist","description":"Belt"},"6":{"gopher":{"id":"2","name":"","description":"","things":null},"id":"6","name":"Feet","description":"Shoes"}}},"4":{"id":"4","name":"Gopher4","description":"The fourth Gopher (#4th)","things":{"11":{"gopher":{"id":"4","name":"","description":"","things":null},"id":"11","name":"Waist","description":"Belt"},"3":{"gopher":{"id":"4","name":"","description":"","things":null},"id":"3","name":"Head","description":"Hat"},"7":{"gopher":{"id":"4","name":"","description":"","things":null},"id":"7","name":"Feet","description":"Shoes"}}},"8":{"id":"8","name":"Gopher8","description":"The eighth Gopher (#8th)","things":{"12":{"gopher":{"id":"8","name":"","description":"","things":null},"id":"12","name":"Waist","description":"Belt"},"4":{"gopher":{"id":"8","name":"","description":"","things":null},"id":"4","name":"Head","description":"Hat"},"8":{"gopher":{"id":"8","name":"","description":"","things":null},"id":"8","name":"Feet","description":"Shoes"}}}}`) {
		t.Fatalf("Unexpected output:%s", jsonOutput)
	}

	server.Stop(c, m)
}

func captureList(adapter *pkgclient.Adapter, cli ui.CLI) string {
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	client.List(adapter, cli)
	w.Close()
	os.Stdout = oldStdout
	var buf bytes.Buffer
	io.Copy(&buf, r)

	return buf.String()
}
func captureAdd(adapter *pkgclient.Adapter, cli ui.CLI) string {
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	client.Add(adapter, cli)
	w.Close()
	os.Stdout = oldStdout
	var buf bytes.Buffer
	io.Copy(&buf, r)

	return buf.String()
}

func captureDelete(adapter *pkgclient.Adapter, cli ui.CLI) string {
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	client.Delete(adapter, cli)
	w.Close()
	os.Stdout = oldStdout
	var buf bytes.Buffer
	io.Copy(&buf, r)

	return buf.String()
}

func SetupConfig(c *config.Config) {
	c.Client.OutputMode = "JSON"
	c.Server.ListenPort = "30303"
	// Use a different port than the DEFAULT, then we can parallel tests
	c.Client.BaseURL = "http://localhost:30303"
	// Test cases are run from the package folder containing the source file.
	c.TemplateFolder = "./../../../../config/template/"
	c.VerboseLevel1 = true
	c.VerboseLevel = "1"

	c.Client.AccessKey = "notempty"
	c.Client.SecretKey = "notempty"

	_ = os.RemoveAll(c.Server.CacheFolder)
	_ = os.RemoveAll(c.Client.CacheFolder)

	c.ValidateOrFatal()
}
