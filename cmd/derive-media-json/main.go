// derive-media-json is a command line tool to derive an abbreviated "media.json" file from a
// "contents/posts-(N).html" file as published by the Instagram export tool, circa April 2022. Previous
// Instagram export data bundles (circa October, 2020) used to provide one or more "media-(N).json"
// files that contained machine-readable properties for working with Instagram exports. This tool
// attempts to reconstruct that data derived from HTML markup and outputs the results as JSON to STDOUT.
// For example:
//
//	$> bin/derive-media-json /usr/local/instagram-export/contents/posts_1.html
package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"os"

	"github.com/sfomuseum/go-sfomuseum-instagram/media"
)

func main() {

	flag.Parse()

	ctx := context.Background()

	photos := make([]*media.Photo, 0)

	paths := flag.Args()

	for _, path := range paths {

		posts_r, err := os.Open(path)

		if err != nil {
			log.Fatalf("Failed to open %s, %v", path, err)
		}

		defer posts_r.Close()

		photos, err = media.DerivePhotosFromReader(ctx, posts_r, photos)

		if err != nil {
			log.Fatalf("Failed to parse photos for %s, %v", path, err)
		}
	}

	archive := media.Archive{
		Photos: photos,
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent(" ", " ")
	enc.Encode(archive)
}
