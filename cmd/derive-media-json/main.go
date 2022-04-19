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
	is_timestamp := false

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

			if n.Type == html.ElementNode && n.Data == "td" {

				first := n.FirstChild

				if first.Data == "Creation Timestamp" {
					is_timestamp = true
				} else if is_timestamp {

					taken = first.Data

					p := &Post{
						MediaId: media_id,
						Path:    path,
						Taken:   taken,
					}

					posts = append(posts, p)

					media_id = ""
					path = ""
					taken = ""

					is_timestamp = false

				} else {
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

func parseComments(ctx context.Context, r io.Reader) (map[string]string, error) {

	doc, err := html.Parse(r)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse HTML, %w", err)
	}

	comments := make(map[string]string)

	is_comment := false
	is_timestamp := false
	is_owner := false

	text := ""
	timestamp := ""
	owner := ""

	var f func(*html.Node)

	f = func(n *html.Node) {

		// fmt.Println(n.Data)

		if n.Type == html.ElementNode && n.Data == "td" {

			first := n.FirstChild

			if first == nil {
				// pass
			} else if first.Data == "Comment" {
				is_comment = true
			} else if first.Data == "Comment creation time" {
				is_timestamp = true
			} else if first.Data == "Media owner" {
				is_owner = true
			} else if is_comment {

				// better error checking; is nil?
				text = first.FirstChild.Data
				is_comment = false

			} else if is_timestamp {

				timestamp = first.Data
				is_timestamp = false

			} else if is_owner {

				// better error checking; is nil?
				owner = first.FirstChild.Data
				is_owner = false

				if owner == "sfomuseum" {

					_, exists := comments[timestamp]

					if exists {
						log.Println("ERP")
					} else {
						comments[timestamp] = text
					}
				}

				text = ""
				timestamp = ""
				owner = ""

			} else {
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	f(doc)

	return comments, nil
}

func main() {

	posts := flag.String("posts", "", "...")
	comments := flag.String("comments", "", "...")

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

	comments_r, err := os.Open(*comments)

	if err != nil {
		log.Fatalf("Failed to open %s, %v", *comments, err)
	}

	defer comments_r.Close()

	c, err := parseComments(ctx, comments_r)

	if err != nil {
		log.Fatalf("Failed to parse comments, %v", err)
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent(" ", " ")
	enc.Encode(p)

	enc.Encode(c)
}
