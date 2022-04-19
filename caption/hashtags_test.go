package caption

import (
	"strings"
	"testing"
)

func TestDeriveHashTagsFromCaption(t *testing.T) {

	tests := map[string][]string{
		"This is a test #test":             []string{"test"},
		"This has no tags":                 []string{},
		"Here is a #complex test #testing": []string{"complex", "testing"},
		"Here is #another test":            []string{"another"},
		"Wishing everyone the happiest, healthiest, and most joyful of New Years! May your dreams take flight in 2022.\n\n#HappyNewYear #2022 #PacificSouthwestAirlines\n\n📸:\nPacific Southwest Airlines (PSA), holiday greeting\nGift of the William Hough Collection\n2006.010.705": []string{"HappyNewYear", "2022", "PacificSouthwestAirlines"},
	}

	for caption, expected := range tests {

		tags, err := DeriveHashTagsFromCaption(caption)

		if err != nil {
			t.Fatalf("Failed to derive hash tags from '%s', %v", caption, err)
		}

		tags_count := len(tags)
		expected_count := len(expected)

		if tags_count != expected_count {
			t.Fatalf("Invalid tag count. Expected %d, got %d", expected_count, tags_count)
		}

		str_tags := strings.Join(tags, "")
		str_expected := strings.Join(expected, "")

		if str_tags != str_expected {
			t.Fatalf("Invalid stringification. Expected '%s', got '%s'", str_expected, str_tags)
		}
	}
}
