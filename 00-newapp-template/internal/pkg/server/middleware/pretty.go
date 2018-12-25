package middleware

import (
	"bytes"
	"net/http"
	"os/exec"
	"strings"
)

// PrettyResponseCtx runs for every route, sets the response to JSON for all responses and unpacks AccessKey&SecretKey
func PrettyResponseCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w = NewPrettyPrint(w)
		next.ServeHTTP(w, r)
	})
}

// NewPrettyPrint checks for jq
func NewPrettyPrint(w http.ResponseWriter) (p *prettyPrintJSON) {
	p = new(prettyPrintJSON)
	p.w = w

	jq, err := exec.LookPath("jq")
	if err == nil {
		p.jq = jq
	}

	return
}

// prettyPrintJSON holds a reference to the ResponseWrite and where 'jq' exec is
type prettyPrintJSON struct {
	w  http.ResponseWriter
	jq string
}

// Write is called, and we rewrite if jq is installed in exec path
func (r *prettyPrintJSON) Write(bb []byte) (int, error) {
	if r.jq == "" {
		return r.w.Write(bb)
	}

	var pretty bytes.Buffer
	raw := bb
	cmd := exec.Command(r.jq, ".")
	cmd.Stdin = strings.NewReader(string(raw))
	cmd.Stdout = &pretty
	err := cmd.Run()
	if err == nil {
		bb = []byte(pretty.String())
	}

	return r.w.Write(bb)
}

// Header overrides Header from ResponseWriter
func (r *prettyPrintJSON) Header() http.Header {
	return r.w.Header()
}

// WriteHeader overrides ResponseWriter
func (r *prettyPrintJSON) WriteHeader(statusCode int) {
	r.w.WriteHeader(statusCode)
}
