package adapter

// Gopher is used by the Client Library and converted from ACME JSONs
type Gopher struct {
	ID          string           `json:"id"` // We prefer strings instead of json.Numbers
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Things      map[string]Thing `json:"things"`
}

// Thing is used by the Client Library and converted from ACME JSONs
type Thing struct {
	Gopher      Gopher `json:"gopher"`
	ID          string `json:"id"` // We prefer strings instead of json.Numbers
	Name        string `json:"name"`
	Description string `json:"description"`
}
