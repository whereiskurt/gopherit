package acme

import (
	"00-newapp-template/pkg/cache"
	"bytes"
	"encoding/json"
	"fmt"
	"gopkg.in/matryer/try.v1"
	"log"
	"strings"
	"sync"
	"text/template"
)

// ServiceMap defines all the endpoints provided by the ACME service
var ServiceMap = map[EndPoint]ServiceTransport{
	EndPointServices.Gophers: {
		URL:           "/gophers",
		CacheFilename: "Gophers.json",
		MethodTemplate: map[httpMethodType]MethodTemplate{
			HTTP.Get: {},
			HTTP.Put: {`{"name": "{{.Name}}", "description":"{{.Description}}"}`},
		},
	},
	EndPointServices.Gopher: {
		URL:           "/gopher/{{.GopherID}}",
		CacheFilename: "gopher/{{.GopherID}}/Gopher.json",
		MethodTemplate: map[httpMethodType]MethodTemplate{
			HTTP.Get:    {},
			HTTP.Delete: {},
			HTTP.Post:   {`{"name": "{{.Name}}", "description":"{{.Description}}"}`},
		},
	},
	EndPointServices.Things: {
		URL:           "/gopher/{{.GopherID}}/things",
		CacheFilename: "gopher/{{.GopherID}}/Things.json",
		MethodTemplate: map[httpMethodType]MethodTemplate{
			HTTP.Get: {},
			HTTP.Put: {`{"name": "{{.Name}}", "description":"{{.Description}}"}`},
		},
	},
	EndPointServices.Thing: {
		URL:           "/gopher/{{.GopherID}}/thing/{{.ThingID}}",
		CacheFilename: "gopher/{{.GopherID}}/thing/{{.ThingID}}/Thing.json",
		MethodTemplate: map[httpMethodType]MethodTemplate{
			HTTP.Get:    {},
			HTTP.Delete: {},
			HTTP.Post:   {`{"name": "{{.Name}}", "description":"{{.Description}}"}`},
		},
	},
}

// ServiceTransport describes a URL endpoint that can be called ACME. Depending on the HTTP method (GET/POST/DELETE)
// we will render the appropriate MethodTemplate
type ServiceTransport struct {
	URL            string
	CacheFilename  string
	MethodTemplate map[httpMethodType]MethodTemplate
}

// MethodTemplate for each GET/PUT/POST/DELETE that is called this template is rendered
// For POST it is put in the BODY, for GET it is added after "?" on the URL, f
type MethodTemplate struct {
	Template string
}

// DefaultRetryIntervals values in here we control the re-try of the Service
var DefaultRetryIntervals = []int{0, 500, 500, 500, 500, 1000, 1000, 1000, 1000, 1000, 3000}

// Service exposes ACME services by converting the JSON results to to Go []structures
type Service struct {
	BaseURL        string // Put in front of every transport call
	SecretKey      string // ACME Secret Keys for API Access (provided by ACME)
	AccessKey      string //             Access Key for API access (provided by ACME)
	RetryIntervals []int  // When a call to a transport fails, this will control the retrying.
	DiskCache      *cache.Disk
	Worker         *sync.WaitGroup // Used by Go routines to control workers (TODO)
}

// NewService is configured to call ACME services with the BaseURL and credentials.
// BaseURL is ofter set to localhost for Unit Testing
func NewService(base string, secret string, access string) (s Service) {
	s.BaseURL = strings.TrimSuffix(base, "/")
	s.SecretKey = secret
	s.AccessKey = access
	s.RetryIntervals = DefaultRetryIntervals
	s.Worker = new(sync.WaitGroup)
	return
}

// EnableCache will create a new Disk Cache for all request.
func (s *Service) EnableCache(cacheFolder string, cryptoKey string) {
	var useCrypto = false
	if cryptoKey != "" {
		useCrypto = true
	}
	s.DiskCache = cache.NewDisk(cacheFolder, cryptoKey, useCrypto)
	return
}

// GetGophers uses a Transport to make GET HTTP call against ACME "GetGophers"
// If the Service RetryIntervals list is populated the calls will retry on Transport errors.
func (s *Service) GetGophers() (gophers []Gopher) {
	tErr := try.Do(func(attempt int) (shouldRetry bool, err error) {
		body, err := s.get(EndPointServices.Gophers, nil)
		if err != nil {
			log.Printf("failed getting gophers: error:%s", err)
			shouldRetry = s.sleepBeforeRetry(attempt)
			return
		}
		// Take the Transport results and convert to []struts
		err = json.Unmarshal(body, &gophers)
		if err != nil {
			shouldRetry = s.sleepBeforeRetry(attempt)
			return
		}

		return
	})
	if tErr != nil {
		log.Printf("failed to GET gophers: %+v", tErr)
	}

	return
}

// GetThings uses a Transport to make a GET HTTP call against ACME "GetThings".
// If the Service RetryIntervals list is populated the calls will retry on Transport errors.
func (s *Service) GetThings(gopherID string) (things []Thing) {
	tErr := try.Do(func(attempt int) (shouldRetry bool, err error) {
		body, err := s.get(EndPointServices.Things, map[string]string{"GopherID": gopherID})
		if err != nil {
			shouldRetry = s.sleepBeforeRetry(attempt)
			return
		}

		err = json.Unmarshal(body, &things)
		if err != nil {
			shouldRetry = s.sleepBeforeRetry(attempt)
			return
		}
		return
	})

	if tErr != nil {
		log.Printf("failed to GET things: %+v", tErr)
	}
	return
}

// DeleteGopher uses a Transport to make a DELETE HTTP call against ACME "DeleteGophers"
// If the Service RetryIntervals list is populated the calls will retry on Transport errors.
func (s *Service) DeleteGopher(gopherID string) (gophers []Gopher) {
	tErr := try.Do(func(attempt int) (shouldRetry bool, err error) {
		body, err := s.delete(EndPointServices.Gopher, map[string]string{"GopherID": gopherID})
		if err != nil {
			log.Printf("failed to DELETE Gopher: %+v", err)
		}
		err = json.Unmarshal(body, &gophers)
		if err != nil {
			log.Printf("failed to unmarshal non-deleted gophers: %+v", err)
		}
		return
	})

	if tErr != nil {
		log.Printf("failed to DELETE gopher: %+v", tErr)
	}
	return
}

// DeleteThing uses a Transport to make a DELETE HTTP call against ACME "DeleteGophers"
// If the Service RetryIntervals list is populated the calls will retry on Transport errors.
func (s *Service) DeleteThing(gopherID string, thingID string) (things []Thing) {
	p := make(map[string]string)
	p["ThingID"] = thingID
	p["GopherID"] = gopherID
	body, err := s.delete(EndPointServices.Thing, p)
	if err != nil {
		log.Printf("failed to DELETE thing: %+v", err)
	}
	err = json.Unmarshal(body, &things)
	if err != nil {
		log.Printf("failed to unmarshal non-deleted things: %+v", err)
	}

	return
}

func ToURL(baseURL string, name EndPoint, p map[string]string) (string, error) {
	sMap, hasMethod := ServiceMap[name]
	if !hasMethod {
		return "", fmt.Errorf("invalid name '%s' for URL lookup", name)
	}

	if p == nil {
		p = make(map[string]string)
	}
	p["BaseURL"] = baseURL

	// Append the BaseURL to the URL
	url := fmt.Sprintf("%s%s", baseURL, sMap.URL)

	return ToTemplate(name, p, url)
}

func ToCacheFilename(name EndPoint, p map[string]string) (string, error) {
	sMap, hasMethod := ServiceMap[name]
	if !hasMethod {
		return "", fmt.Errorf("invalid name '%s' for cache filename lookup", name)
	}
	return ToTemplate(name, p, sMap.CacheFilename)
}

func ToJSON(name EndPoint, method httpMethodType, p map[string]string) (string, error) {
	sMap, hasMethod := ServiceMap[name]
	if !hasMethod {
		return "", fmt.Errorf("invalid method '%s' for name '%s'", method, name)
	}

	mMap, hasTemplate := sMap.MethodTemplate[method]
	if !hasTemplate {
		return "", fmt.Errorf("invalid template for method '%s' for name '%s'", method, name)
	}

	tmpl := mMap.Template
	return ToTemplate(name, p, tmpl)
}

func ToTemplate(name EndPoint, data map[string]string, tmpl string) (string, error) {
	var rawURL bytes.Buffer
	t, terr := template.New(fmt.Sprintf("%s", name)).Parse(tmpl)
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
