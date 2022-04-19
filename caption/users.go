package caption

import (
	"unicode"
)

// DeriveUserNamesFromCaption returns the list of user name (@{USERNAME}) contained in 'body'
func DeriveUserNamesFromCaption(body string) ([]string, error) {

	names := make([]string, 0)
	n := ""

	capture := false

	for _, r := range body {

		if string(r) == "@" {
			capture = true
			continue
		}

		if capture && unicode.IsSpace(r) {
			capture = false

			if n != "" {
				names = append(names, n)
				n = ""
			}
		}

		if capture {
			n = n + string(r)
		}
	}

	if n != "" {
		names = append(names, n)
	}

	return names, nil
}
