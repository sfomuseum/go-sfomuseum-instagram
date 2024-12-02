package media

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// TIME_FORMAT is a time.Parse compatible string representing the manner in which Instagram datetime strings are encoded.
// I mean that's the idea anyway. IG seems to change to format they use between exports... Basically what that means is
// you should the `ParseTime` function (in this package) to parse IG dates.
const TIME_FORMAT string = "Jan 2, 2006 3:04 PM"

// TIME_FORMAT_UPPER is a time.Parse compatible string representing the manner in which Instagram datetimes encoded with upper-case AM/PM times.
const TIME_FORMAT_UPPER string = "Jan 2, 2006 3:04 PM"

// TIME_FORMAT_UPPER_COMMA is a time.Parse compatible string representing the manner in which Instagram datetimes encoded with upper-case AM/PM times
// with dates and times separated by a comma.
const TIME_FORMAT_UPPER_COMMA string = "Jan 2, 2006, 3:04 PM"

// TIME_FORMAT_LOWER is a time.Parse compatible string representing the manner in which Instagram datetimes encoded with lower-case AM/PM times.
const TIME_FORMAT_LOWER string = "Jan 2, 2006 3:04 pm"

// TIME_FORMAT_LOWER_COMMA is a time.Parse compatible string representing the manner in which Instagram datetimes encoded with lower-case AM/PM times
// with dates and times separated by a comma.
const TIME_FORMAT_LOWER_COMMA string = "Jan 2, 2006, 3:04 pm"

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

	t, err := ParseTime(str_created)

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

func ParseTime(str_time string) (time.Time, error) {

	layouts := []string{
		TIME_FORMAT_UPPER,
		TIME_FORMAT_LOWER,
		TIME_FORMAT_UPPER_COMMA,
		TIME_FORMAT_LOWER_COMMA,
	}

	var t time.Time
	var err error

	for _, layout := range layouts {

		t, err = time.Parse(layout, str_time)

		if err != nil {
			slog.Debug("Failed to parse time", "time", str_time, "layout", layout, "error", err)
			continue
		}

		break
	}

	if err != nil {
		t := new(time.Time)
		return *t, err
	}

	return t, nil
}
