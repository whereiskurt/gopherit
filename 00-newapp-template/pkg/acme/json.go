package acme

import (
	"encoding/json"
)

// Gopher is returned by the ACME services and defined by ACME documentation
type Gopher struct {
	ID          json.Number `json:"gopher_id"` // ACME uses json.Numbers
	Name        string      `json:"name"`
	Description string      `json:"description"`
}

// Thing is returned by the ACME services and defined by ACME documentation
type Thing struct {
	ID          json.Number `json:"thing_id"` // ACME uses json.Numbers
	GopherID    json.Number `json:"gopher_id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
}
