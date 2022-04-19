package main

import (
	"context"
	"crypto/sha1"
	"encoding/json"
	"flag"
	"fmt"
	"golang.org/x/net/html"
	"io"
	"log"
	"os"
	"path/filepath"
)

type Caption struct {
	Excerpt string `json:"excerpt"`
}

type Post struct {
	MediaId string  `json:"media_id"`
	Path    string  `json:"path"`
	Taken   string  `json:"taken"`
	Caption Caption `json:"caption"`
}

type Media struct {
	Posts []Post
}

func parsePosts(ctx context.Context, r io.Reader) ([]*Post, error) {

	doc, err := html.Parse(r)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse HTML, %w", err)
	}

	posts := make([]*Post, 0)

	is_post := false

	var media_id string
	var path string
	var taken string

	var f func(*html.Node)

	f = func(n *html.Node) {

		// fmt.Println(n.Data)

		if n.Type == html.ElementNode && n.Data == "div" {

			for _, a := range n.Attr {

				switch a.Key {
				case "class":

					if a.Val == "_a6-p" {
						is_post = true
					}

				default:
					// pass
				}
			}
		}

		if is_post {

			if path != "" {

				p := &Post{
					MediaId: media_id,
					Path:    path,
					Taken:   taken,
				}

				posts = append(posts, p)

				media_id = ""
				path = ""
				taken = ""
			}
		}

		if is_post {

			if n.Type == html.ElementNode && n.Data == "img" {

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

	posts := flag.String("posts", "", "...")
	// comments := flag.String("comments", "", "...")

	flag.Parse()

	ctx := context.Background()

	posts_r, err := os.Open(*posts)

	if err != nil {
		log.Fatalf("Failed to open %s, %v", *posts, err)
	}

	defer posts_r.Close()

	p, err := parsePosts(ctx, posts_r)

	if err != nil {
		log.Fatalf("Failed to parse posts, %v", err)
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent(" ", " ")
	enc.Encode(p)
}
