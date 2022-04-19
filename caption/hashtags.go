package caption

import (
	"unicode"
)

// DeriveHashTagsFromCaption returns the list of hash tags contained in 'body'
func DeriveHashTagsFromCaption(body string) ([]string, error) {

	tags := make([]string, 0)
	t := ""

	capture := false

	for _, r := range body {

		if string(r) == "#" {
			capture = true
			continue
		}

		if capture && unicode.IsSpace(r) {
			capture = false

			if t != "" {
				tags = append(tags, t)
				t = ""
			}
		}

		if capture {
			t = t + string(r)
		}
	}

	if t != "" {
		tags = append(tags, t)
	}

	return tags, nil
}
