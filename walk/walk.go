package walk

import (
	"context"
	"encoding/json"
	"github.com/sfomuseum/go-sfomuseum-instagram"
	"io"
	"sync"
)

type WalkOptions struct {
	MediaChannel chan []byte
	ErrorChannel chan error
	DoneChannel  chan bool
}

type WalkMediaCallbackFunc func(ctx context.Context, media []byte) error

func WalkMediaWithCallback(ctx context.Context, media_fh io.Reader, cb WalkMediaCallbackFunc) error {

	err_ch := make(chan error)
	media_ch := make(chan []byte)
	done_ch := make(chan bool)

	walk_opts := &WalkOptions{
		DoneChannel:  done_ch,
		ErrorChannel: err_ch,
		MediaChannel: media_ch,
	}

	go WalkMedia(ctx, walk_opts, media_fh)

	working := true
	wg := new(sync.WaitGroup)

	for {
		select {
		case <-done_ch:
			working = false
		case err := <-err_ch:
			return err
		case body := <-media_ch:

			wg.Add(1)

			go func(body []byte) {

				defer wg.Done()

				err := cb(ctx, body)

				if err != nil {
					err_ch <- err
				}

			}(body)

		}

		if !working {
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

	var archive instagram.Archive

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

		opts.MediaChannel <- ph_body
	}

}
