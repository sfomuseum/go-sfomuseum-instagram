package document

import (
	"context"
	"errors"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"regexp"
	"strings"
)

var re_hashtag *regexp.Regexp
var re_newlines *regexp.Regexp

func init() {

	re_hashtag = regexp.MustCompile(`.*((#([^#\s]+)\s?)+)$`)
	re_newlines = regexp.MustCompile(`(\.?\n\.)+$`)
}

type Caption struct {
	Body     string   `json:"body"`
	Hashtags []string `json:"hashtags"`
}

func ExpandCaption(ctx context.Context, body []byte) ([]byte, error) {

	caption_rsp := gjson.GetBytes(body, "caption")

	if !caption_rsp.Exists() {
		return nil, errors.New("Missing caption")
	}

	str_caption := caption_rsp.String()

	caption, err := ParseCaption(ctx, str_caption)

	if err != nil {
		return nil, err
	}

	body, err = sjson.SetBytes(body, "caption", caption)

	if err != nil {
		return nil, err
	}

	return body, nil
}

func ParseCaption(ctx context.Context, body string) (*Caption, error) {

	m := re_hashtag.FindStringSubmatch(body)

	str_hashtags := ""
	hashtags := make([]string, 0)

	if len(m) > 0 {

		str_hashtags = m[0]

		trim_hashtags := strings.TrimSpace(str_hashtags)

		for _, tag := range strings.Split(trim_hashtags, " ") {

			tag = strings.Replace(tag, "#", "", 1)
			tag = strings.TrimSpace(tag)

			if tag != "" {
				hashtags = append(hashtags, tag)
			}
		}
	}

	if str_hashtags != "" {
		body = strings.Replace(body, str_hashtags, "", 1)
	}

	body = strings.TrimSpace(body)

	body = re_newlines.ReplaceAllString(body, "")

	caption := &Caption{
		Body:     body,
		Hashtags: hashtags,
	}

	return caption, nil
}
