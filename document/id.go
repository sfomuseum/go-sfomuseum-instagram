package document

import (
	"context"
	"errors"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"path/filepath"
	"strings"
)

// AppendMediaIDFromPath will derive a media ID from the `path` JSON property
// in 'body' and use that value to append a `media_id` JSON property to 'body'.
// Note: this function uses code specific to pre- April, 2022 exports and does
// not (yet) match what we do in media/json.go. This will need to be resolved.
func AppendMediaIDFromPath(ctx context.Context, body []byte) ([]byte, error) {

	path_rsp := gjson.GetBytes(body, "path")

	if !path_rsp.Exists() {
		return nil, errors.New("Missing path property")
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
// Note: this function uses code specific to pre- April, 2022 exports and does
// not (yet) match what we do in media/json.go. This will need to be resolved.
// Given that Instagram appears to have changed the filenames for media we will need
// to do ... something. Image or file hashes perhaps?
func DeriveMediaIDFromPath(path string) string {

	fname := filepath.Base(path)
	ext := filepath.Ext(path)

	id := strings.Replace(fname, ext, "", 1)
	return id
}
