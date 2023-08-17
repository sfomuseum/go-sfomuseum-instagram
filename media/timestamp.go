package media

import (
	"context"
	"fmt"
	"time"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// TIME_FORMAT is a time.Parse compatible string representing the manner in which Instagram datetime strings are encoded.
const TIME_FORMAT string = "Jan 2, 2006, 3:04 PM"

// AppendTakenCallback is a user-defined callback function to be applied to `time.Time` instances in the
// `AppendTakenAtTimestampWithCallback` method after an initial datetime string has been parsed but before
// a Unix timestamp is appeneded (to an Instagram post).
type AppendTakenCallback func(time.Time) (time.Time, error)

// AppendTakenAtTimestamp will look for a `taken_at` JSON property in 'body' and use
// its value to derive a Unix timestamp which will be used as the value of a new
// `taken` JSON property (which is appended to 'body').
func AppendTakenAtTimestamp(ctx context.Context, body []byte) ([]byte, error) {
	return AppendTakenAtTimestampWithCallback(ctx, body, nil)
}

func AppendTakenAtTimestampWithCallback(ctx context.Context, body []byte, cb AppendTakenCallback) ([]byte, error) {

	created_rsp := gjson.GetBytes(body, "taken_at")

	if !created_rsp.Exists() {
		return nil, fmt.Errorf("Missing taken_at property")
	}

	str_created := created_rsp.String()

	t, err := time.Parse(TIME_FORMAT, str_created)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse taken value (%s), %w", str_created, err)
	}

	if cb != nil {

		t, err = cb(t)

		if err != nil {
			return nil, fmt.Errorf("Custom time parsing callback failed, %w", err)
		}
	}

	body, err = sjson.SetBytes(body, "taken", t.Unix())

	if err != nil {
		return nil, fmt.Errorf("Failed to assign taken property, %w", err)
	}

	return body, nil
}
