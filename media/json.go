package media

import (
	"context"
	"crypto/sha1"
	"fmt"
	"golang.org/x/net/html"
	"io"
	"net/url"
	"path/filepath"
	"time"
)

// DerivePhotosFromReader will derive zero or more Instagram photos from the body of 'r' appending
// each to 'photos'.
func DerivePhotosFromReader(ctx context.Context, r io.Reader, photos []*Photo) ([]*Photo, error) {

	doc, err := html.Parse(r)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse HTML, %w", err)
	}

	var media_id string
	var path string
	var taken string
	var caption string

	var f func(*html.Node)

	f = func(n *html.Node) {

		if n.Type == html.ElementNode {

			if n.Data == "div" {

				is_caption := false
				is_taken := false

				for _, a := range n.Attr {

					switch a.Key {
					case "class":

						if a.Val == "_3-95 _2pim _a6-h _a6-i" {
							is_caption = true
						}

						if a.Val == "_3-94 _a6-o" {
							is_taken = true
						}

					default:
						// pass
					}
				}

				if is_caption {
					caption = n.FirstChild.Data
					is_caption = false
				}

				if is_taken {

					taken = n.FirstChild.Data
					is_taken = false

					taken_at := ""

					t, err := time.Parse("Jan 2, 2006, 3:04 PM", taken)

					if err == nil {
						taken_at = t.Format(time.RFC3339)
					}

					// log.Println(taken ,taken_at)

					if path != "" {

						fname := filepath.Base(path)
						data := []byte(fname)
						media_id = fmt.Sprintf("%x", sha1.Sum(data))

						p := &Photo{
							Path:    path,
							TakenAt: taken_at,
							Caption: caption,
							MediaId: media_id,
						}

						photos = append(photos, p)
					}

					path = ""
					media_id = ""
					caption = ""
					taken = ""
				}

			} else if n.Data == "a" {

				for _, a := range n.Attr {

					switch a.Key {
					case "href":

						if filepath.Ext(a.Val) == ".mp4" {

							u, err := url.Parse(a.Val)

							if err == nil {
								path = u.Path
							}

						}

					default:
						// pass
					}

				}

			} else if n.Data == "img" {

				for _, a := range n.Attr {

					switch a.Key {
					case "src":
						path = a.Val
					default:
						// pass
					}
				}

			} else {
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	f(doc)

	return photos, nil
}
