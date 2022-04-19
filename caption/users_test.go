package caption

import (
	"strings"
	"testing"
)

func TestDeriveUserNamesFromCaption(t *testing.T) {

	tests := map[string][]string{
		"This is a test @test":             []string{"test"},
		"This has no names":                []string{},
		"Here is a @complex test @testing": []string{"complex", "testing"},
		"Here is @another test":            []string{"another"},
		"Jane Chow, a Los Angeles-based filmmaker from Hong Kong, produces this award-winning film about a lonely teenager who tries to help her parents keep their restaurant afloat during the COVID-19 pandemic in Los Angeles Chinatown.\n\nSee “Sorry for the Inconvenience” by Jane Chow in the Video Arts Gallery, located pre-security, in the International Terminal with daily operating hours of 5:00am to 10:00pm or online at: https://bit.ly/3E3DE7J\n\n#VideoArts #VideoArtsSFO #JaneChow #SorryForTheInconvenience \n@janialls": []string{"janialls"},
	}

	for caption, expected := range tests {

		names, err := DeriveUserNamesFromCaption(caption)

		if err != nil {
			t.Fatalf("Failed to derive hash names from '%s', %v", caption, err)
		}

		names_count := len(names)
		expected_count := len(expected)

		if names_count != expected_count {
			t.Fatalf("Invalid tag count. Expected %d, got %d", expected_count, names_count)
		}

		str_names := strings.Join(names, "")
		str_expected := strings.Join(expected, "")

		if str_names != str_expected {
			t.Fatalf("Invalid stringification. Expected '%s', got '%s'", str_expected, str_names)
		}
	}
}
