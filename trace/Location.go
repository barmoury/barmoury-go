package trace

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	State     string  `json:"state"`
	Country   string  `json:"country"`
	Address   string  `json:"address"`
}
