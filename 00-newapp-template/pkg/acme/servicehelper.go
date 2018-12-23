package acme

import (
	"bytes"
	"fmt"
	"text/template"
	"time"
)

// ServiceTransport describes a URL endpoint that can be called ACME. Depending on the HTTP method (GET/POST/DELETE)
// we will render the appropriate MethodTemplate
type ServiceTransport struct {
	URL            string
	CacheFilename  string
	MethodTemplate map[string]MethodTemplate
}

// MethodTemplate for each GET/PUT/POST/DELETE that is called this template is rendered
// For POST it is put in the BODY, for GET it is added after "?" on the URL, f
type MethodTemplate struct {
	Template string
}

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

	// If we have a DiskCache it means we will write out responses to disk.
	if s.DiskCache != nil {
		// We have initialized a cache then write to it.
		filename, err := s.toCacheFilename(name, p)
		if err != nil {
			return nil, err
		}
		filename = fmt.Sprintf("%s/%s", s.DiskCache.CacheFolder , filename)
		err = s.DiskCache.Store(filename, body)
		if err != nil {
			return nil, err
		}
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
	json, err := s.toJSON(name, "POST", p)
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

func (s *Service) toURL(name string, p map[string]string) (string, error) {
	sMap, hasMethod := serviceMap[name]
	if !hasMethod {
		return "", fmt.Errorf("invalid name '%s' for URL lookup", name)
	}

	if p == nil {
		p = make(map[string]string)
	}
	p["BaseURL"] = s.BaseURL

	// Append the BaseURL to the URL
	url := fmt.Sprintf("%s%s", s.BaseURL, sMap.URL)

	return s.toTemplate(name, p, url)
}

func (s *Service) toCacheFilename(name string, p map[string]string) (string, error) {
	sMap, hasMethod := serviceMap[name]
	if !hasMethod {
		return "", fmt.Errorf("invalid name '%s' for cache filename lookup", name)
	}
	return s.toTemplate(name, p, sMap.CacheFilename)
}

func (s *Service) toJSON(name string, method string, p map[string]string) (string, error) {
	sMap, hasMethod := serviceMap[name]
	if !hasMethod {
		return "", fmt.Errorf("invalid method '%s' for name '%s'", method, name)
	}

	mMap, hasTemplate := sMap.MethodTemplate[method]
	if !hasTemplate {
		return "", fmt.Errorf("invalid template for method '%s' for name '%s'", method, name)
	}

	tmpl := mMap.Template
	return s.toTemplate(name, p, tmpl)
}

func (s *Service) toTemplate(name string, data map[string]string, tmpl string) (string, error) {
	var rawURL bytes.Buffer
	t, terr := template.New(name).Parse(tmpl)
	if terr != nil {
		err := fmt.Errorf("error: failed to parse template for %s: %v", name, terr)
		return "", err
	}
	err := t.Execute(&rawURL, data)
	if err != nil {
		return "", err
	}

	url := rawURL.String()

	return url, nil
}
