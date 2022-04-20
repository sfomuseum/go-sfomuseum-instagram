package media

import (
	"bytes"
	"context"
	"fmt"
	"github.com/sfomuseum/go-sfomuseum-instagram/hash"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"gocloud.dev/blob"
	"io"
	"path/filepath"
)

// AppendHashes will append a variety of hashes to 'body' derived from its contents.
// For images it will append a "file_hash" (SHA-256) and a "perceptual_hash" JSON properties.
// For images it will append a "file_hash" (SHA-256) JSON property.
// Media file URIs are derived from the "path" JSON property in 'body'. These URIs are expected
// to be relative and resolvable in 'bucket'.
func AppendHashes(ctx context.Context, body []byte, bucket *blob.Bucket) ([]byte, error) {

	updates := make(map[string]string)

	path_rsp := gjson.GetBytes(body, "path")

	if !path_rsp.Exists() {
		return nil, fmt.Errorf("Missing path property")
	}

	rel_path := path_rsp.String()

	media_fh, err := bucket.NewReader(ctx, rel_path, nil)

	if err != nil {
		return nil, fmt.Errorf("Failed to open %s, %w", rel_path, err)
	}

	defer media_fh.Close()

	media_body, err := io.ReadAll(media_fh)

	if err != nil {
		return nil, fmt.Errorf("Failed to read %s, %w", rel_path, err)
	}

	media_r := bytes.NewReader(media_body)

	file_hash, err := hash.FileHash(media_r)

	if err != nil {
		return nil, fmt.Errorf("Failed to generate file hash for %s, %w", rel_path, err)
	}

	updates["file_hash"] = file_hash

	if filepath.Ext(rel_path) != ".mp4" {

		media_r.Seek(0, 0)

		p_hash, err := hash.PerceptualHash(media_r)

		if err != nil {
			return nil, fmt.Errorf("Failed to generate file perceptual for %s, %w", rel_path, err)
		}

		updates["perceptual_hash"] = p_hash
	}

	for k, v := range updates {

		body, err = sjson.SetBytes(body, k, v)

		if err != nil {
			return nil, fmt.Errorf("Failed to assign %s (%s), %w", k, v, err)
		}
	}

	return body, nil
}
