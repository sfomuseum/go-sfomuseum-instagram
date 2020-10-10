package instagram

type Archive struct {
	Photos []Photo `json:"photos"`
}

type Photo struct {
	Caption  string `json:"caption"`
	TakenAt  string `json:"taken_at"` // to do time.Time parsing
	Location string `json:"location"`
	Path     string `json:"path"`
}
