package acme

import (
	"bytes"
	"fmt"
	"text/template"
	"time"
)

func (s *Service) sleepBeforeRetry(attempt int) (shouldReRun bool) {
	if attempt < len(s.RetryIntervals) {
		time.Sleep(time.Duration(s.RetryIntervals[attempt]) * time.Millisecond)
		shouldReRun = true
	}
	return
}

func (s *Service) get(name string, p map[string]string) ([]byte, error) {
	url, err := s.toURL(name, p)
	if err != nil {
		return nil, err
	}
	t := NewTransport(s)
	body, err := t.get(url)
	if err != nil {
		return nil, err
	}
	return body, err
}
func (s *Service) delete(name string, p map[string]string) ([]byte, error) {
	url, err := s.toURL(name, p)
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
func (s *Service) update(name string, p map[string]string) ([]byte, error) {
	url, err := s.toURL(name, p)
	if err != nil {
		return nil, err
	}
	json, err := s.toJSON(name, p)
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

func (s *Service) toURL(name string, p map[string]string) (url string, err error) {
	if p == nil {
		p = make(map[string]string)
	}
	p["BaseURL"] = s.BaseURL
	return s.toTemplate(name, p, urlServiceTmpl)
}
func (s *Service) toJSON(name string, p map[string]string) (url string, err error) {
	return s.toTemplate(name, p, jsonBodyTmpl)
}

func (s *Service) toTemplate(name string, p map[string]string, tmap map[string]string) (url string, err error) {
	var rawURL bytes.Buffer
	t, terr := template.New(name).Parse(tmap[name])
	if terr != nil {
		err = fmt.Errorf("error: failed to parse template for %s: %v", name, err)
		return
	}
	err = t.Execute(&rawURL, p)
	if err != nil {
		return
	}

	url = rawURL.String()

	return
}
