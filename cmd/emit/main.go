package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/aaronland/go-json-query"
	"github.com/sfomuseum/go-sfomuseum-instagram/media"
	"github.com/sfomuseum/go-sfomuseum-instagram/walk"
	"github.com/tidwall/pretty"
	_ "gocloud.dev/blob/fileblob"
)

func main() {

	media_uri := flag.String("media-uri", "", "A valid gocloud.dev/blob URI to your Instagram `media.json` file.")

	to_stdout := flag.Bool("stdout", true, "Emit to STDOUT")
	to_devnull := flag.Bool("null", false, "Emit to /dev/null")
	as_json := flag.Bool("json", false, "Emit a JSON list.")
	format_json := flag.Bool("format-json", false, "Format JSON output for each record.")

	append_timestamp := flag.Bool("append-timestamp", false, "Append a `taken` property containing a Unix timestamp derived from the `taken_at` property.")
	append_id := flag.Bool("append-id", false, "Append a `media_id` property derived from the `path` property.")
	append_all := flag.Bool("append-all", false, "Enable all the `-append-` flags.")

	expand_caption := flag.Bool("expand-caption", false, "Parse and replace the string `caption` property with a `media.Caption` struct.")

	var queries query.QueryFlags
	flag.Var(&queries, "query", "One or more {PATH}={REGEXP} parameters for filtering records.")

	valid_modes := strings.Join([]string{query.QUERYSET_MODE_ALL, query.QUERYSET_MODE_ANY}, ", ")
	desc_modes := fmt.Sprintf("Specify how query filtering should be evaluated. Valid modes are: %s", valid_modes)

	query_mode := flag.String("query-mode", query.QUERYSET_MODE_ALL, desc_modes)

	flag.Parse()

	if *append_all {
		*append_timestamp = true
		*append_id = true
	}

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

	media_fh, err := media.Open(ctx, *media_uri)

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

		if *append_timestamp {

			b, err := media.AppendTakenAtTimestamp(ctx, body)

			if err != nil {
				return err
			}

			body = b
		}

		if *append_id {

			b, err := media.AppendMediaIDFromPath(ctx, body)

			if err != nil {
				return err
			}

			body = b
		}

		if *expand_caption {

			b, err := media.ExpandCaption(ctx, body)

			if err != nil {
				return err
			}

			body = b
		}

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

	walk_opts := &walk.WalkWithCallbackOptions{
		Callback: cb,
	}

	if len(queries) > 0 {

		qs := &query.QuerySet{
			Queries: queries,
			Mode:    *query_mode,
		}

		walk_opts.QuerySet = qs
	}

	err = walk.WalkMediaWithCallback(ctx, walk_opts, media_fh)

	if err != nil {
		log.Fatal(err)
	}

	if *as_json {
		wr.Write([]byte("]"))
	}

}
