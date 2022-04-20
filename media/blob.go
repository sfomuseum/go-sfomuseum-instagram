package media

import (
	"context"
	"fmt"
	"gocloud.dev/blob"
	"io"
	"path/filepath"
)

func Open(ctx context.Context, media_uri string) (io.ReadCloser, error) {

	root := filepath.Dir(media_uri)
	fname := filepath.Base(media_uri)

	media_bucket, err := blob.OpenBucket(ctx, root)

	if err != nil {
		return nil, fmt.Errorf("Failed to open bucket (%s), %v", root, err)
	}

	return media_bucket.NewReader(ctx, fname, nil)
}
