package document

import (
	"context"
	"errors"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"time"
)

func AppendTakenAtTimestamp(ctx context.Context, body []byte) ([]byte, error) {

	created_rsp := gjson.GetBytes(body, "taken_at")

	if !created_rsp.Exists() {
		return nil, errors.New("Missing taken_at property")
	}

	str_created := created_rsp.String()

	// "taken_at": "2020-10-07T00:34:36+00:00"

	t_fmt := time.RFC3339
	t, err := time.Parse(t_fmt, str_created)

	if err != nil {
		return nil, err
	}

	body, err = sjson.SetBytes(body, "taken", t.Unix())

	if err != nil {
		return nil, err
	}

	return body, nil
}
