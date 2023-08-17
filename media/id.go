package media

import (
	"context"
	"crypto/sha1"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// DeriveMediaIdFromString returns the SHA-1 value of 'input'. Because Instagram
// doesn't include stable identifiers in its media.json output we need to derive
// one from some element associated with each post/photo. We used to do this using
// the file path but since those changed sometime between October 2020 and April 2022.
// So rather than picking a specific key whose value or formatting may change again
// this method exists to consisitently hash an arbitrary string supplied by consumers
// of this package. None of this should be necessary but until there are stable IDs
// the Instagram exports this is what we get.
func DeriveMediaIdFromString(input string) string {
	data := []byte(input)
	return fmt.Sprintf("%x", sha1.Sum(data))
}

// AppendMediaIDFromPath will derive a media ID from the `path` JSON property
// in 'body' and use that value to append a `media_id` JSON property to 'body'.
// This method is deprecated.
func AppendMediaIDFromPath(ctx context.Context, body []byte) ([]byte, error) {

	path_rsp := gjson.GetBytes(body, "path")

	if !path_rsp.Exists() {
		return nil, fmt.Errorf("Missing path property")
	}

	path := path_rsp.String()

	id := DeriveMediaIDFromPath(path)

	body, err := sjson.SetBytes(body, "media_id", id)

	if err != nil {
		return nil, err
	}

	return body, nil
}

// DeriveMediaIDFromPath will return a (old) media ID derived from 'path'.
// This method is deprecated.
func DeriveMediaIDFromPath(path string) string {

	fname := filepath.Base(path)
	ext := filepath.Ext(path)

	id := strings.Replace(fname, ext, "", 1)
	return id
}
