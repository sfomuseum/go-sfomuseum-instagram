package media

import (
	"context"
	"fmt"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"time"
)

// AppendTakenAtTimestamp will look for a `taken_at` JSON property in 'body' and use
// its value to derive a Unix timestamp which will be used as the value of a new
// `taken` JSON property (which is appended to 'body').
func AppendTakenAtTimestamp(ctx context.Context, body []byte) ([]byte, error) {

	created_rsp := gjson.GetBytes(body, "taken_at")

	if !created_rsp.Exists() {
		return nil, fmt.Errorf("Missing taken_at property")
	}

	str_created := created_rsp.String()

	// "taken_at": "2020-10-07T00:34:36+00:00"

	t_fmt := time.RFC3339
	t, err := time.Parse(t_fmt, str_created)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse taken value (%s), %w", str_created, err)
	}

	body, err = sjson.SetBytes(body, "taken", t.Unix())

	if err != nil {
		return nil, fmt.Errorf("Failed to assign taken property, %w", err)
	}

	return body, nil
}
