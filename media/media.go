// package media provides methods for working with or deriving media-(N).json files.
package media

type Archive struct {
	Photos []*Photo `json:"photos"`
}

type Photo struct {
	Caption  string `json:"caption"`
	TakenAt  string `json:"taken_at"` // to do time.Time parsing
	Location string `json:"location,omitempty"`
	Path     string `json:"path"`
}
