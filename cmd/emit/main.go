package main

import (
	"context"
	"flag"
	"github.com/sfomuseum/go-sfomuseum-instagram/walk"
	"github.com/tidwall/pretty"
	"gocloud.dev/blob"
	_ "gocloud.dev/blob/fileblob"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
)

func main() {

	media_uri := flag.String("media-uri", "", "A valid gocloud.dev/blob URI to your Instagram `media.json` file.")

	to_stdout := flag.Bool("stdout", true, "Emit to STDOUT")
	to_devnull := flag.Bool("null", false, "Emit to /dev/null")
	as_json := flag.Bool("json", false, "Emit a JSON list.")
	format_json := flag.Bool("format-json", false, "Format JSON output for each record.")

	flag.Parse()

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	defer cancel()

	writers := make([]io.Writer, 0)

	if *to_stdout {
		writers = append(writers, os.Stdout)
	}

	if *to_devnull {
		writers = append(writers, ioutil.Discard)
	}

	if len(writers) == 0 {
		log.Fatal("Nothing to write to.")
	}

	wr := io.MultiWriter(writers...)

	root := filepath.Dir(*media_uri)
	fname := filepath.Base(*media_uri)

	media_bucket, err := blob.OpenBucket(ctx, root)

	if err != nil {
		log.Fatalf("Failed to open bucket (%s), %v", root, err)
	}

	defer media_bucket.Close()

	media_fh, err := media_bucket.NewReader(ctx, fname, nil)

	if err != nil {
		log.Fatalf("Failed to open media file, %v", err)
	}

	defer media_fh.Close()

	count := uint32(0)
	mu := new(sync.RWMutex)

	if *as_json {
		wr.Write([]byte("["))
	}

	cb := func(ctx context.Context, body []byte) error {

		mu.Lock()
		defer mu.Unlock()

		new_count := atomic.AddUint32(&count, 1)

		if new_count > 1 {

			if *as_json {
				wr.Write([]byte(","))
			}
		}

		if *as_json && *format_json {
			body = pretty.Pretty(body)
		}

		wr.Write(body)
		wr.Write([]byte("\n"))

		return nil
	}

	err = walk.WalkMediaWithCallback(ctx, media_fh, cb)

	if err != nil {
		log.Fatal(err)
	}

	if *as_json {
		wr.Write([]byte("]"))
	}

}
