package walk

import (
	"context"
	"encoding/json"
	"github.com/sfomuseum/go-sfomuseum-instagram"
	"io"
	"sync"
)

type WalkOptions struct {
	PhotoChannel chan []byte
	ErrorChannel chan error
	DoneChannel  chan bool
}

type WalkPhotosCallbackFunc func(ctx context.Context, photo []byte) error

func WalkPhotosWithCallback(ctx context.Context, photos_fh io.Reader, cb WalkPhotosCallbackFunc) error {

	err_ch := make(chan error)
	photo_ch := make(chan []byte)
	done_ch := make(chan bool)

	walk_opts := &WalkOptions{
		DoneChannel:  done_ch,
		ErrorChannel: err_ch,
		PhotoChannel: photo_ch,
	}

	go WalkPhotos(ctx, walk_opts, photos_fh)

	working := true
	wg := new(sync.WaitGroup)

	for {
		select {
		case <-done_ch:
			working = false
		case err := <-err_ch:
			return err
		case body := <-photo_ch:

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

func WalkPhotos(ctx context.Context, opts *WalkOptions, photos_fh io.Reader) {

	defer func() {
		opts.DoneChannel <- true
	}()

	var archive instagram.Archive

	dec := json.NewDecoder(photos_fh)
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

		opts.PhotoChannel <- ph_body
	}

}
