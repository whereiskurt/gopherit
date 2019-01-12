// +build !release

package config

import (
	"net/http"
)

// TemplateFolder implements the http filesystem, but is overridden when we
// build with tags (go build -tags release) this file won't be built, but
// templates_generate.go will be.
var TemplateFolder http.FileSystem = http.Dir("config/")