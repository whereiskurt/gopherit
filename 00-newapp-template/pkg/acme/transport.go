package acme

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"
)

var tr = &http.Transport{
	MaxIdleConns:    20,
	IdleConnTimeout: 30 * time.Second,
}

// Transport defines the HTTP details for the API call.
type Transport struct {
	BaseURL     string
	AccessKey   string
	SecretKey   string
	WorkerCount int
	ThreadSafe  *sync.Mutex
}

// NewTransport handles the HTTP methods GET/POST/PUT/DELETE
func NewTransport(s *Service) (p Transport) {
	p.BaseURL = s.BaseURL
	p.AccessKey = s.AccessKey
	p.SecretKey = s.SecretKey
	p.ThreadSafe = new(sync.Mutex)
	return
}

// Inserts the AccessKey and SecretKey into the request header.
// AccessKey/SecretKey may be equally lengthed comma separated values that are rotated through each call.
// headerCallCount is thread-safely incremented allowing multiple-requests from multiple-credentials (access/secret keys.)
var headerCallCount int

func (t *Transport) header() string {
	akeys := strings.Split(t.AccessKey, ",")
	skeys := strings.Split(t.SecretKey, ",")

	if len(akeys) != len(skeys) {
		return ""
	}

	// Ensure incremental non-overalapping count
	t.ThreadSafe.Lock()
	headerCallCount = headerCallCount + 1
	mod := headerCallCount % len(akeys)
	t.ThreadSafe.Unlock()

	return fmt.Sprintf("AccessKey=%s;SecretKey=%s", akeys[mod], skeys[mod])
}

func (t *Transport) get(url string) (body []byte, err error) {
	var req *http.Request
	var resp *http.Response

	client := &http.Client{Transport: tr}

	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	req.Header.Add("X-ApiKeys", t.header())

	resp, err = client.Do(req) // <-------HTTPS GET Request!
	if err != nil {
		return
	}

	if resp.StatusCode == 429 {
		err = errors.New("error: we need to slow down")
		return
	}
	if resp.StatusCode == 403 {
		err = errors.New("error: creds no longer authorized")
		return
	}

	if resp.StatusCode != 200 {
		err = fmt.Errorf("error: status code does not appear successful: %d", resp.StatusCode)
		return
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err == nil {
		err = resp.Body.Close()
	}

	return
}
func (t *Transport) post(url string, data string, datatype string) (body []byte, err error) {
	var req *http.Request
	var resp *http.Response

	client := &http.Client{Transport: tr}

	req, err = http.NewRequest("POST", url, bytes.NewBuffer([]byte(data)))
	if err != nil {
		return
	}
	req.Header.Add("X-ApiKeys", t.header())
	req.Header.Set("Content-Type", datatype)

	resp, err = client.Do(req) // <-------HTTPS GET Request!
	if err != nil {
		return
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err == nil {
		err = resp.Body.Close()
	}
	return
}
func (t *Transport) put(url string) (body []byte, err error) {
	var req *http.Request
	var resp *http.Response

	client := &http.Client{Transport: tr}

	req, err = http.NewRequest("put", url, nil)
	if err != nil {
		return
	}

	req.Header.Add("X-ApiKeys", t.header())

	resp, err = client.Do(req)
	if err != nil {
		return
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err == nil {
		err = resp.Body.Close()
	}
	return
}
func (t *Transport) delete(url string) (body []byte, err error) {
	var req *http.Request
	var resp *http.Response

	client := &http.Client{Transport: tr}

	req, err = http.NewRequest("DELETE", url, nil)
	if err != nil {
		return
	}

	req.Header.Add("X-ApiKeys", t.header())

	resp, err = client.Do(req)
	if err != nil {
		return
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err == nil {
		err = resp.Body.Close()
	}
	return
}
