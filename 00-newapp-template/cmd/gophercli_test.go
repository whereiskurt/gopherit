package main

import (
	"os"
	"testing"
)

func TestGopherCLI(t *testing.T) {
	// NOTE: Don't use " or ', this isn't a shell and we don't need them - they get passed literally!
	os.Args = []string{
		"gopherit", "client", "list",
		"--mode=json",
		"--configFolder=./../config/",
		"--configFile=default.test.gophercli",
		"--templateFolder=./../config/template",
	}
	main()
}
