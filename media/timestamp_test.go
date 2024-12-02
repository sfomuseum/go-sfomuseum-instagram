package media

import (
	"testing"
)

func TestParseTime(t *testing.T) {

	str_times := []string{
		"Mar 10, 2023, 1:25 PM",
		"Mar 10, 2023, 1:25 am",
		"Mar 10, 2023 1:25 PM",
		"Mar 10, 2023 1:25 am",
	}

	for _, str_t := range str_times {

		_, err := ParseTime(str_t)

		if err != nil {
			t.Fatalf("Failed to parse '%s', %v", str_t, err)
		}
	}
}
