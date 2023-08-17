package media

import (
	"context"
	"fmt"
	"io"
	"path/filepath"

	"gocloud.dev/blob"
)

// Open is a convenience method to create a new `blob.Bucket` instance derived from the
// root directory in 'media_uri' and then using that bucket to return a new `io.ReadCloser`
// instance for 'media_uri'.
func Open(ctx context.Context, media_uri string) (io.ReadCloser, error) {

	root := filepath.Dir(media_uri)
	fname := filepath.Base(media_uri)

	media_bucket, err := blob.OpenBucket(ctx, root)

	if err != nil {
		return nil, fmt.Errorf("Failed to open bucket (%s), %v", root, err)
	}

	return media_bucket.NewReader(ctx, fname, nil)
}
