package acme

import (
	"encoding/json"
	"gopkg.in/matryer/try.v1"
	"log"
	"strings"
	"sync"
)

var serviceMap = map[string]ServiceTransport{
	"Gophers": {
		URL: "{{.BaseURL}}/gophers",
		MethodTemplate: map[string]MethodTemplate{
			"GET": {},
		},
	},
	"Gopher": {
		URL: "{{.BaseURL}}/gopher/{{.GopherID}}",
		MethodTemplate: map[string]MethodTemplate{
			"GET":    {""},
			"DELETE": {""},
			"PUT":    {`{"name": "{{.Name}}", "description":"{{.Description}}"}`},
			"POST":   {`{"name": "{{.Name}}", "description":"{{.Description}}"}`},
		},
	},
	"Things": {
		URL: "{{.BaseURL}}/gopher/{{.GopherID}}/things",
		MethodTemplate: map[string]MethodTemplate{
			"GET": {""},
		},
	},
	"Thing": {
		URL: "{{.BaseURL}}/gopher/{{.GopherID}}/thing/{{.ThingID}}",
		MethodTemplate: map[string]MethodTemplate{
			"GET":    {""},
			"DELETE": {""},
			"PUT":    {`{"name": "{{.Name}}", "description":"{{.Description}}"}`},
			"POST":   {`{"name": "{{.Name}}", "description":"{{.Description}}"}`},
		},
	},
}

// DefaultRetryIntervals values in here we control the re-try of the Service
var DefaultRetryIntervals = []int{0, 500, 500, 500, 500, 1000, 1000, 1000, 1000, 1000, 3000}


// Service exposes ACME services by converting the JSON results to to Go []structures
type Service struct {
	BaseURL        string          // Put in front of every
	SecretKey      string          // ACME Secret Keys for API Access (provided by ACME)
	AccessKey      string          //             Access Key for API access (provided by ACME)
	RetryIntervals []int           // When a call to a transport fails, this will control the retrying.
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

// GetGophers uses a Transport to make GET HTTP call against ACME "GetGophers"
// If the Service RetryIntervals list is populated the calls will retry on Transport errors.
func (s *Service) GetGophers() (gophers []Gopher) {
	tErr := try.Do(func(attempt int) (shouldRetry bool, err error) {
		body, err := s.get("Gophers", nil)
		if err != nil {
			shouldRetry = s.sleepBeforeRetry(attempt)
			return
		}
		// Take the Transport results and convert to []struts
		err = json.Unmarshal(body, &gophers)
		if err != nil {
			shouldRetry = s.sleepBeforeRetry(attempt)
			return
		}

		// TODO: Write to cache!
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
		body, err := s.get("Things", map[string]string{"GopherID": gopherID})
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
		body, err := s.delete("Gopher", map[string]string{"GopherID": gopherID})
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
	body, err := s.delete("Thing", p)
	if err != nil {
		log.Printf("failed to DELETE thing: %+v", err)
	}
	err = json.Unmarshal(body, &things)
	if err != nil {
		log.Printf("failed to unmarshal non-deleted things: %+v", err)
	}

	return
}
