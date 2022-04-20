package media

// Type Archive is a struct representing the structure of an Instagram media.json file.
type Archive struct {
	// Photos is the list of photos (posts) in an archive.
	Photos []*Photo `json:"photos"`
}

// Type Photo is a struct containing data associated with an Instagram post. The name `Photo` reflects the
// naming conventions of the (old) Instagram media.json files.
type Photo struct {
	// Caption is the caption associated with the post
	Caption string `json:"caption"`
	// Taken is the datetime string when the post was published
	TakenAt  string `json:"taken_at"`
	Location string `json:"location,omitempty"`
	// Path is the relative URI for the media element associated with the post
	Path    string `json:"path"`
	MediaId string `json:"media_id,omitempty"`
}
