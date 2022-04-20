package media

import (
	"context"
	"crypto/sha1"
	"fmt"
	_ "github.com/sfomuseum/go-sfomuseum-instagram/caption"
	"golang.org/x/net/html"
	"io"
	"log"
	"net/url"
	"path/filepath"
	"time"
)

// type Caption is a struct containing data associated with the caption for an Instragram psot
type Caption struct {
	// Excerpt is the body of the caption
	Excerpt string `json:"excerpt,omitempty"`
	// Body is the body of the caption
	Body string `json:"body"`
	// HashTags is the list of hash tags contained in the body of the caption
	HashTags []string `json:"hashtags,omitempty"`
	// Users is the list of user names contained in the body of the caption
	Users []string `json:"users,omitempty"`
}

// type Post is a struct containing data associated with an Instagram post
type Post struct {
	// MediaId is the SHA-1 hash of the basename for the path of the media element associated with the post
	MediaId string `json:"media_id"`
	// Path is the relative URI for the media element associated with the post
	Path string `json:"path"`
	// Taken is the datetime string when the post was published
	TakenAt string `json:"taken_at"`
	// Caption is the caption associated with the post
	Caption *Caption `json:"caption"`
}

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
	var body string

	// var tags []string
	// var users []string

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

					body = n.FirstChild.Data

					/*
						t, err := caption.DeriveHashTagsFromCaption(body)

						if err == nil {
							tags = t
						}

						u, err := caption.DeriveUserNamesFromCaption(body)

						if err == nil {
							users = u
						}
					*/

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

						/*
							c := &Caption{
								Body:     body,
								HashTags: tags,
								Users:    users,
							}
						*/

						p := &Photo{
							Path: path,
							// MediaId: media_id,
							TakenAt: taken_at,
							Caption: body,
						}

						photos = append(photos, p)
					}

					log.Println(media_id)

					path = ""
					media_id = ""
					body = ""
					taken = ""

					// tags = []string{}
					// users = []string{}

				}

			} else if n.Data == "a" {

				for _, a := range n.Attr {

					switch a.Key {
					case "href":

						if filepath.Ext(a.Val) == ".mp4" {

							u, err := url.Parse(a.Val)

							if err == nil {

								path = u.Path
								fname := filepath.Base(path)

								data := []byte(fname)
								media_id = fmt.Sprintf("%x", sha1.Sum(data))
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
						fname := filepath.Base(path)

						data := []byte(fname)
						media_id = fmt.Sprintf("%x", sha1.Sum(data))

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
