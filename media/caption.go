package media

import (
	"context"
	"errors"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"gopkg.in/neurosnap/sentences.v1/english"
	"regexp"
	"strings"
)

var re_hashtag *regexp.Regexp
var re_separator *regexp.Regexp
var re_newlines *regexp.Regexp

func init() {

	re_hashtag = regexp.MustCompile(`.*(((?:#|@)([^#@\s]+)\s?)+)$`)
	re_separator = regexp.MustCompile(`(\.?\n\.)+$`)
	re_newlines = regexp.MustCompile(`\n`)
}

type Caption struct {
	Body     string   `json:"body"`
	Excerpt  string   `json:"excerpt"`
	Hashtags []string `json:"hashtags"`
	Users    []string `json:"users"`
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
	users := make([]string, 0)

	if len(m) > 0 {

		str_hashtags = m[0]

		trim_hashtags := strings.TrimSpace(str_hashtags)

		for _, tag := range strings.Split(trim_hashtags, " ") {

			if strings.HasPrefix(tag, "#") {
				tag = strings.Replace(tag, "#", "", 1)
				tag = strings.TrimSpace(tag)

				if tag != "" {
					hashtags = append(hashtags, tag)
				}

			} else if strings.HasPrefix(tag, "@") {

				tag = strings.Replace(tag, "@", "", 1)
				tag = strings.TrimSpace(tag)

				if tag != "" {
					users = append(users, tag)
				}

			} else {
				// this should never happen
			}
		}
	}

	if str_hashtags != "" {
		body = strings.Replace(body, str_hashtags, "", 1)
	}

	body = strings.TrimSpace(body)

	body = re_separator.ReplaceAllString(body, "")
	body = re_newlines.ReplaceAllString(body, " ")

	caption := &Caption{
		Body:     body,
		Hashtags: hashtags,
		Users:    users,
	}

	tokenizer, err := english.NewSentenceTokenizer(nil)

	if err != nil {
		return nil, err
	}

	sentences := tokenizer.Tokenize(body)

	if len(sentences) >= 1 {
		caption.Excerpt = sentences[0].Text
	}

	return caption, nil
}
