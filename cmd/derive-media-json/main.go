// derive-media-json is a command line tool to derive an abbreviated "media.json" file from a
// "contents/posts-(N).html" file as published by the Instagram export tool, circa April 2022. Previous
// Instagram export data bundles (circa October, 2020) used to provide one or more "media-(N).json"
// files that contained machine-readable properties for working with Instagram exports. This tool
// attempts to reconstruct that data derived from HTML markup and outputs the results as JSON to STDOUT.
// For example:
// 	$> bin/derive-media-json /usr/local/instagram-export/contents/posts_1.html
package main

import (
	"context"
	"crypto/sha1"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/sfomuseum/go-sfomuseum-instagram/caption"
	"golang.org/x/net/html"
	"io"
	"log"
	"os"
	"path/filepath"
)

// type Caption is a struct containing data associated with the caption for an Instragram psot
type Caption struct {
	// Excerpt is the body of the caption
	Excerpt string `json:"excerpt,omitempty"`
	// Body is the body of the caption
	Body     string   `json:"body"`
	HashTags []string `json:"hashtags,omitempty"`
	Users    []string `json:"users,omitempty"`
}

// type Post is a struct containing data associated with an Instagram post
type Post struct {
	// MediaId is the SHA-1 hash of the basename for the path of the media element associated with the post
	MediaId string `json:"media_id"`
	// Path is the relative URI for the media element associated with the post
	Path string `json:"path"`
	// Taken is the datetime string when the post was published
	Taken string `json:"taken"`
	// Caption is the caption associated with the post
	Caption *Caption `json:"caption"`
}

// DerivePostsFromReader will derive zero or more Instagram posts from the body of 'r' appending
// each to 'posts'.
func DerivePostsFromReader(ctx context.Context, r io.Reader, posts []*Post) ([]*Post, error) {

	doc, err := html.Parse(r)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse HTML, %w", err)
	}

	var media_id string
	var path string
	var taken string
	var body string

	var tags []string
	var users []string

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

					t, err := caption.DeriveHashTagsFromCaption(body)

					if err == nil {
						tags = t
					}

					u, err := caption.DeriveUserNamesFromCaption(body)

					if err == nil {
						users = u
					}

					is_caption = false
				}

				if is_taken {

					taken = n.FirstChild.Data
					is_taken = false

					if path != "" {

						log.Println(users)
						c := &Caption{
							Body:     body,
							HashTags: tags,
							Users:    users,
						}

						p := &Post{
							Path:    path,
							MediaId: media_id,
							Taken:   taken,
							Caption: c,
						}

						posts = append(posts, p)
					}

					path = ""
					media_id = ""
					body = ""
					taken = ""

					tags = []string{}
					users = []string{}

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

	return posts, nil
}

func main() {

	flag.Parse()

	ctx := context.Background()

	posts := make([]*Post, 0)

	paths := flag.Args()

	for _, path := range paths {

		posts_r, err := os.Open(path)

		if err != nil {
			log.Fatalf("Failed to open %s, %v", path, err)
		}

		defer posts_r.Close()

		posts, err = DerivePostsFromReader(ctx, posts_r, posts)

		if err != nil {
			log.Fatalf("Failed to parse posts for %s, %v", path, err)
		}
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent(" ", " ")
	enc.Encode(posts)
}
