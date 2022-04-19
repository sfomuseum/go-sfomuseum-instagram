package instagram

import (
	"context"
	"fmt"
	"gocloud.dev/blob"
	"io"
	"path/filepath"
)

type Archive struct {
	Photos []*Photo `json:"photos"`
}

type Photo struct {
	Caption  string `json:"caption"`
	TakenAt  string `json:"taken_at"` // to do time.Time parsing
	Location string `json:"location,omitempty"`
	Path     string `json:"path"`
}

func OpenMedia(ctx context.Context, media_uri string) (io.ReadCloser, error) {

	root := filepath.Dir(media_uri)
	fname := filepath.Base(media_uri)

	media_bucket, err := blob.OpenBucket(ctx, root)

	if err != nil {
		return nil, fmt.Errorf("Failed to open bucket (%s), %v", root, err)
	}

	return media_bucket.NewReader(ctx, fname, nil)
}
