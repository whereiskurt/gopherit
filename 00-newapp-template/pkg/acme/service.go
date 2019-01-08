package acme

import (
	"00-newapp-template/pkg/cache"
	"00-newapp-template/pkg/metrics"
	"bytes"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"gopkg.in/matryer/try.v1"
	"strings"
	"sync"
	"text/template"
	"time"
)

// DefaultRetryIntervals values in here we control the re-try of the Service
var DefaultRetryIntervals = []int{0, 500, 500, 500, 500, 1000, 1000, 1000, 1000, 1000, 3000}

var EndPoints = endPointTypes{
	Gophers: EndPointType("Gophers"),
	Gopher:  EndPointType("Gopher"),
	Things:  EndPointType("Things"),
	Thing:   EndPointType("Thing"),
}

// ServiceMap defines all the endpoints provided by the ACME service
var ServiceMap = map[EndPointType]ServiceTransport{
	EndPoints.Gophers: {
		URL:           "/gophers",
		CacheFilename: "Gophers.json",
		MethodTemplate: map[httpMethodType]MethodTemplate{
			HTTP.Put: {`{{.GopherJSON}}`},
		},
	},
	EndPoints.Gopher: {
		URL:           "/gopher/{{.GopherID}}",
		CacheFilename: "gopher/{{.GopherID}}/Gopher.json",
		MethodTemplate: map[httpMethodType]MethodTemplate{
			HTTP.Post: {`{{.GopherJSON}}`},
		},
	},
	EndPoints.Things: {
		URL:           "/gopher/{{.GopherID}}/things",
		CacheFilename: "gopher/{{.GopherID}}/Things.json",
		MethodTemplate: map[httpMethodType]MethodTemplate{
			HTTP.Put: {`{{.ThingJSON}}`},
		},
	},
	EndPoints.Thing: {
		URL:           "/gopher/{{.GopherID}}/thing/{{.ThingID}}",
		CacheFilename: "gopher/{{.GopherID}}/thing/{{.ThingID}}/Thing.json",
		MethodTemplate: map[httpMethodType]MethodTemplate{
			HTTP.Post: {`{{.ThingJSON}}`},
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
// For POST it is Put in the BODY, for GET it is added after "?" on the URL, f
type MethodTemplate struct {
	Template string
}

// Service exposes ACME services by converting the JSON results to to Go []structures
type Service struct {
	BaseURL        string // Put in front of every transport call
	SecretKey      string // ACME Secret Keys for API Access (provided by ACME)
	AccessKey      string //             Access Key for API access (provided by ACME)
	RetryIntervals []int  // When a call to a transport fails, this will control the retrying.
	DiskCache      *cache.Disk
	Worker         *sync.WaitGroup // Used by Go routines to control workers (TODO)
	Log            *log.Logger
	Metrics        *metrics.Metrics
}

type EndPointType string
type endPointTypes struct {
	Gophers EndPointType
	Gopher  EndPointType
	Things  EndPointType
	Thing   EndPointType
}

func (c EndPointType) String() string {
	return "pkg.acme.endpoints." + string(c)
}

// NewService is configured to call ACME services with the BaseURL and credentials.
// BaseURL is ofter set to localhost for Unit Testing
func NewService(base string, secret string, access string) (s Service) {
	s.BaseURL = strings.TrimSuffix(base, "/")
	s.SecretKey = secret
	s.AccessKey = access
	s.RetryIntervals = DefaultRetryIntervals
	s.Worker = new(sync.WaitGroup)
	s.Log = new(log.Logger)
	return
}

func (s *Service) EnableMetrics(metrics *metrics.Metrics) {
	s.Metrics = metrics
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

func (s *Service) SetLogger(log *log.Logger) {
	s.Log = log
}

// GetGophers uses a Transport to make GET HTTP call against ACME "GetGophers"
// If the Service RetryIntervals list is populated the calls will retry on Transport errors.
func (s *Service) GetGophers() (gophers []Gopher) {

	tErr := try.Do(func(attempt int) (shouldRetry bool, err error) {
		body, status, err := s.get(EndPoints.Gophers, nil)

		if s.Metrics != nil {
			s.Metrics.TransportInc(metrics.EndPoints.Gophers, metrics.Methods.Transport.Get, status)
		}

		if err != nil {
			s.Log.Warnf("failed getting gophers: error:%s: %d", err, status)
			shouldRetry = s.sleepBeforeRetry(attempt)
			return
		}
		// Take the Transport results and convert to []struts
		err = json.Unmarshal(body, &gophers)
		if err != nil {
			s.Log.Warnf("failed to unmarshal gophers: %s: ", err)
			shouldRetry = s.sleepBeforeRetry(attempt)
			return
		}

		return
	})
	if tErr != nil {
		s.Log.Warnf("failed to GET gophers: %+v", tErr)
	}

	return
}

// GetThings uses a Transport to make a GET HTTP call against ACME "GetThings".
// If the Service RetryIntervals list is populated the calls will retry on Transport errors.
func (s *Service) GetThings(gopherID string) (things []Thing) {
	tErr := try.Do(func(attempt int) (shouldRetry bool, err error) {
		body, status, err := s.get(EndPoints.Things, map[string]string{"GopherID": gopherID})
		if s.Metrics != nil {
			s.Metrics.TransportInc(metrics.EndPoints.Things, metrics.Methods.Transport.Get, status)
		}
		if err != nil {
			s.Log.Infof("failed to retrieve gopherID: %s : http status: %d: %s", gopherID, status, err)
			shouldRetry = s.sleepBeforeRetry(attempt)
			return
		}

		err = json.Unmarshal(body, &things)
		if err != nil {
			s.Log.Infof("failed to unmarshal gopherID: %s: %s", gopherID, err)
			shouldRetry = s.sleepBeforeRetry(attempt)
			return
		}
		return
	})

	if tErr != nil {
		s.Log.Warnf("failed to GET things: %+v", tErr)
	}
	return
}

func (s *Service) AddGopher(g Gopher) Gopher {
	var gopher Gopher

	gjson, err := json.Marshal(g)
	if err != nil {
		return gopher
	}

	_ = try.Do(func(attempt int) (shouldRetry bool, err error) {
		body, status, err := s.add(EndPoints.Gophers, map[string]string{
			"GopherID":   string(g.ID),
			"GopherJSON": string(gjson),
		})
		if s.Metrics != nil {
			s.Metrics.TransportInc(metrics.EndPoints.Gopher, metrics.Methods.Transport.Put, status)
		}
		if status == 403 {
			// FORBIDDEN so don't keep retrying.
			return false, err
		}
		if err != nil {
			s.Log.Warnf("failed to ADD Gopher: %+v", err)
			shouldRetry = s.sleepBeforeRetry(attempt)
			return shouldRetry, err
		}

		err = json.Unmarshal(body, &gopher)
		if err != nil {
			s.Log.Warnf("failed to unmarshal ADDED gopher: %+v", err)
			shouldRetry = s.sleepBeforeRetry(attempt)
		}

		return shouldRetry, err
	})

	return gopher
}

func (s *Service) UpdateGopher(g Gopher) Gopher {
	var gopher Gopher

	gjson, err := json.Marshal(g)

	if err != nil {
		return gopher
	}
	_ = try.Do(func(attempt int) (shouldRetry bool, err error) {
		body, status, err := s.update(EndPoints.Gopher, map[string]string{
			"GopherID":   string(g.ID),
			"GopherJSON": string(gjson),
		})
		if s.Metrics != nil {
			s.Metrics.TransportInc(metrics.EndPoints.Gopher, metrics.Methods.Transport.Post, status)
		}
		if status == 403 {
			// FORBIDDEN so don't keep retrying.
			return false, err
		}
		if err != nil {
			s.Log.Warnf("failed to UPDATE Gopher: %+v", err)
			shouldRetry = s.sleepBeforeRetry(attempt)
			return shouldRetry, err
		}

		err = json.Unmarshal(body, &gopher)
		if err != nil {
			s.Log.Warnf("failed to unmarshal updated gopher: %+v", err)
			shouldRetry = s.sleepBeforeRetry(attempt)
		}

		return shouldRetry, err
	})

	return gopher
}

func (s *Service) UpdateThing(t Thing) Thing {
	var thing Thing

	tjson, err := json.Marshal(t)

	if err != nil {
		return thing
	}
	_ = try.Do(func(attempt int) (shouldRetry bool, err error) {
		body, status, err := s.update(EndPoints.Thing, map[string]string{
			"GopherID":  string(t.GopherID),
			"ThingID":   string(t.ID),
			"ThingJSON": string(tjson),
		})
		if s.Metrics != nil {
			s.Metrics.TransportInc(metrics.EndPoints.Thing, metrics.Methods.Transport.Post, status)
		}
		if status == 403 {
			// FORBIDDEN so don't keep retrying.
			return false, err
		}
		if err != nil {
			s.Log.Warnf("failed to UPDATE Thing: %+v", err)
			shouldRetry = s.sleepBeforeRetry(attempt)
			return shouldRetry, err
		}

		err = json.Unmarshal(body, &thing)
		if err != nil {
			s.Log.Warnf("failed to unmarshal updated thing: %+v", err)
			shouldRetry = s.sleepBeforeRetry(attempt)
		}

		return shouldRetry, err
	})

	return thing
}

// DeleteGopher uses a Transport to make a DELETE HTTP call against ACME "DeleteGophers"
// If the Service RetryIntervals list is populated the calls will retry on Transport errors.
func (s *Service) DeleteGopher(gopherID string) []Gopher {
	var gophers []Gopher

	_ = try.Do(func(attempt int) (shouldRetry bool, err error) {
		body, status, err := s.delete(EndPoints.Gopher, map[string]string{
			"GopherID": gopherID,
		})
		if s.Metrics != nil {
			s.Metrics.TransportInc(metrics.EndPoints.Gopher, metrics.Methods.Transport.Delete, status)
		}

		if status == 403 {
			// FORBIDDEN so don't keep retrying.
			return false, err
		}

		if err != nil {
			s.Log.Warnf("failed to DELETE Gopher: %+v", err)
			shouldRetry = s.sleepBeforeRetry(attempt)
			return shouldRetry, err
		}

		err = json.Unmarshal(body, &gophers)
		if err != nil {
			s.Log.Warnf("failed to unmarshal non-deleted gophers: %+v", err)
			shouldRetry = s.sleepBeforeRetry(attempt)
		}

		return shouldRetry, err
	})

	return gophers
}

// DeleteThing uses a Transport to make a DELETE HTTP call against ACME "DeleteGophers"
// If the Service RetryIntervals list is populated the calls will retry on Transport errors.
func (s *Service) DeleteThing(gopherID string, thingID string) []Thing {
	var things []Thing

	_ = try.Do(func(attempt int) (shouldRetry bool, err error) {
		body, status, err := s.delete(EndPoints.Thing, map[string]string{
			"GopherID": gopherID,
			"ThingID":  thingID,
		})

		if s.Metrics != nil {
			s.Metrics.TransportInc(metrics.EndPoints.Thing, metrics.Methods.Transport.Delete, status)
		}

		if status == 403 {
			// FORBIDDEN so don't keep retrying.
			return false, err
		}

		if err != nil {
			s.Log.Warnf("failed to DELETE thing: %+v: %d", err, status)
			shouldRetry = s.sleepBeforeRetry(attempt)
			return shouldRetry, err
		}

		err = json.Unmarshal(body, &things)
		if err != nil {
			s.Log.Warnf("failed to unmarshal non-deleted things: %+v", err)
		}

		return shouldRetry, err
	})

	return things
}

func ToURL(baseURL string, name EndPointType, p map[string]string) (string, error) {
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

func ToCacheFilename(name EndPointType, p map[string]string) (string, error) {
	sMap, ok := ServiceMap[name]
	if !ok {
		return "", fmt.Errorf("invalid name '%s' for cache filename lookup", name)
	}
	return ToTemplate(name, p, sMap.CacheFilename)
}

func ToJSON(name EndPointType, method httpMethodType, p map[string]string) (string, error) {
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

func ToTemplate(name EndPointType, data map[string]string, tmpl string) (string, error) {
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

func (s *Service) sleepBeforeRetry(attempt int) (shouldReRun bool) {
	if attempt < len(s.RetryIntervals) {
		time.Sleep(time.Duration(s.RetryIntervals[attempt]) * time.Millisecond)
		shouldReRun = true
	}
	return
}

func (s *Service) get(endPoint EndPointType, p map[string]string) ([]byte, int, error) {

	url, err := ToURL(s.BaseURL, endPoint, p)
	if err != nil {
		return nil, 0, err
	}

	t := NewTransport(s)
	body, status, err := t.Get(url)

	if err != nil {
		return nil, status, err
	}

	// If we have a DiskCache it means we will write out responses to disk.
	if s.DiskCache != nil {
		// We have initialized a cache then write to it.
		filename, err := ToCacheFilename(endPoint, p)
		if err != nil {
			return nil, status, err
		}

		err = s.DiskCache.Store(filename, body)
		if err != nil {
			return nil, status, err
		}
	}

	return body, status, err
}
func (s *Service) delete(endPoint EndPointType, p map[string]string) ([]byte, int, error) {
	url, err := ToURL(s.BaseURL, endPoint, p)
	if err != nil {
		return nil, 0, err
	}
	t := NewTransport(s)
	body, status, err := t.Delete(url)
	if err != nil {
		return nil, status, err
	}

	return body, status, err
}
func (s *Service) update(endPoint EndPointType, p map[string]string) ([]byte, int, error) {
	url, err := ToURL(s.BaseURL, endPoint, p)
	if err != nil {
		return nil, 0, err
	}
	j, err := ToJSON(endPoint, HTTP.Post, p)
	if err != nil {
		return nil, 0, err
	}

	t := NewTransport(s)
	body, status, err := t.Post(url, j, "application/json")
	if err != nil {
		return nil, status, err
	}

	return body, status, err
}
func (s *Service) add(endPoint EndPointType, p map[string]string) ([]byte, int, error) {
	url, err := ToURL(s.BaseURL, endPoint, p)
	if err != nil {
		return nil, 0, err
	}
	j, err := ToJSON(endPoint, HTTP.Put, p)
	if err != nil {
		return nil, 0, err
	}

	t := NewTransport(s)
	body, status, err := t.Put(url, j, "application/json")
	if err != nil {
		return nil, status, err
	}

	return body, status, err
}
