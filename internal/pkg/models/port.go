// package contains domain models.
package models

// Port represents the port object in the JSON file
type Port struct {
	PortID      string        `json:"port_id"`
	Name        string        `json:"name"`
	City        string        `json:"city"`
	Province    string        `json:"province"`
	Country     string        `json:"country"`
	Alias       []string      `json:"alias"`
	Regions     []interface{} `json:"regions"`
	Coordinates []float64     `json:"coordinates"`
	Timezone    string        `json:"timezone"`
	Unlocs      []string      `json:"unlocs"`
	Code        string        `json:"code"`
}
