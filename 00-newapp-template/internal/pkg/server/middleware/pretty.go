package middleware

import (
	"bytes"
	"net/http"
	"os/exec"
	"strings"
)

// PrettyCtx runs for every route, sets the response to JSON for all responses and unpacks AccessKey&SecretKey
func PrettyCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w = NewPrettyPrint(w)
		next.ServeHTTP(w, r)
	})
}

type PrettyPrintJSON struct {
	w http.ResponseWriter
	jq string
}

func (r *PrettyPrintJSON) Header() http.Header {
	return r.w.Header()
}

func (r *PrettyPrintJSON) WriteHeader(statusCode int) {
	r.w.WriteHeader(statusCode)
}

func NewPrettyPrint(w http.ResponseWriter) (p *PrettyPrintJSON) {
	p = new(PrettyPrintJSON)
	p.w = w

	jq, err := exec.LookPath("jq")
	if err == nil {
		p.jq = jq
	}

	return
}

func (r *PrettyPrintJSON) Write(bb []byte) (int, error) {
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
