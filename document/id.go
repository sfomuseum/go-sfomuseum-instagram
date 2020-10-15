package document

import (
	"context"
	"errors"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"path/filepath"
	"strings"
)

func AppendMediaIDFromPath(ctx context.Context, body []byte) ([]byte, error) {

	path_rsp := gjson.GetBytes(body, "path")

	if !path_rsp.Exists() {
		return nil, errors.New("Missing path property")
	}

	path := path_rsp.String()

	fname := filepath.Base(path)
	ext := filepath.Ext(path)

	id := strings.Replace(fname, ext, "", 1)

	body, err := sjson.SetBytes(body, "media_id", id)

	if err != nil {
		return nil, err
	}

	return body, nil
}
