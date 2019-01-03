package acme

import (
	"time"
)

type EndPoint string

var EndPointServices = endPointTypes{
	Gophers: EndPoint("Gophers"),
	Gopher:  EndPoint("Gopher"),
	Things:  EndPoint("Things"),
	Thing:   EndPoint("Thing"),
}

type endPointTypes struct {
	Gophers EndPoint
	Gopher  EndPoint
	Things  EndPoint
	Thing   EndPoint
}

func (c EndPoint) String() string {
	return "pkg.acme.endpoint." + string(c)
}

func (s *Service) sleepBeforeRetry(attempt int) (shouldReRun bool) {
	if attempt < len(s.RetryIntervals) {
		time.Sleep(time.Duration(s.RetryIntervals[attempt]) * time.Millisecond)
		shouldReRun = true
	}
	return
}

//TODO: Add 'put' aka 'add'
func (s *Service) get(endPoint EndPoint, p map[string]string) ([]byte, error) {

	url, err := ToURL(s.BaseURL, endPoint, p)
	if err != nil {
		return nil, err
	}
	t := NewTransport(s)
	body, err := t.get(url)
	if err != nil {
		return nil, err
	}

	// If we have a DiskCache it means we will write out responses to disk.
	if s.DiskCache != nil {
		// We have initialized a cache then write to it.
		filename, err := ToCacheFilename(endPoint, p)
		if err != nil {
			return nil, err
		}

		err = s.DiskCache.Store(filename, body)
		if err != nil {
			return nil, err
		}
	}

	return body, err
}

func (s *Service) delete(endPoint EndPoint, p map[string]string) ([]byte, error) {
	url, err := ToURL(s.BaseURL, endPoint, p)
	if err != nil {
		return nil, err
	}
	t := NewTransport(s)
	body, err := t.delete(url)
	if err != nil {
		return nil, err
	}

	return body, err
}
func (s *Service) update(endPoint EndPoint, p map[string]string) ([]byte, error) {
	url, err := ToURL(s.BaseURL, endPoint, p)
	if err != nil {
		return nil, err
	}
	json, err := ToJSON(endPoint, HTTP.Post, p)
	if err != nil {
		return nil, err
	}

	t := NewTransport(s)
	body, err := t.post(url, json, "application/json")
	if err != nil {
		return nil, err
	}

	return body, err
}
