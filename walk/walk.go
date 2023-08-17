package walk

import (
	"context"
	"encoding/json"
	"io"
	_ "log"
	"sync"

	"github.com/aaronland/go-json-query"
	"github.com/sfomuseum/go-sfomuseum-instagram/media"
)

type WalkOptions struct {
	MediaChannel chan []byte
	ErrorChannel chan error
	DoneChannel  chan bool
	QuerySet     *query.QuerySet
}

type WalkWithCallbackOptions struct {
	Callback WalkMediaCallbackFunc
	QuerySet *query.QuerySet
}

type WalkMediaCallbackFunc func(ctx context.Context, media []byte) error

func WalkMediaWithCallback(ctx context.Context, opts *WalkWithCallbackOptions, media_fh io.Reader) error {

	err_ch := make(chan error)
	media_ch := make(chan []byte)
	done_ch := make(chan bool)

	walk_opts := &WalkOptions{
		DoneChannel:  done_ch,
		ErrorChannel: err_ch,
		MediaChannel: media_ch,
		QuerySet:     opts.QuerySet,
	}

	go WalkMedia(ctx, walk_opts, media_fh)

	walking := true
	wg := new(sync.WaitGroup)

	for {
		select {
		case <-done_ch:
			walking = false
		case err := <-err_ch:
			return err
		case body := <-media_ch:

			wg.Add(1)

			go func(body []byte) {

				defer wg.Done()

				err := opts.Callback(ctx, body)

				if err != nil {
					err_ch <- err
				}

			}(body)

		}

		if !walking {
			break
		}
	}

	wg.Wait()
	return nil
}

func WalkMedia(ctx context.Context, opts *WalkOptions, media_fh io.Reader) {

	defer func() {
		opts.DoneChannel <- true
	}()

	var archive media.Archive

	dec := json.NewDecoder(media_fh)
	err := dec.Decode(&archive)

	if err != nil {
		opts.ErrorChannel <- err
		return
	}

	for _, ph := range archive.Photos {

		select {
		case <-ctx.Done():
			break
		default:
			// pass
		}

		ph_body, err := json.Marshal(ph)

		if err != nil {
			opts.ErrorChannel <- err
			continue
		}

		if opts.QuerySet != nil {

			matches, err := query.Matches(ctx, opts.QuerySet, ph_body)

			if err != nil {

				opts.ErrorChannel <- err
				continue
			}

			if !matches {
				continue
			}
		}

		opts.MediaChannel <- ph_body
	}

}
