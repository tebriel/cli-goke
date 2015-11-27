package songs

import (
	"github.com/tebriel/cli-goke/webutils"
	"golang.org/x/net/html"
	// "io/ioutil"
	"io"
	"net/http"
	"strings"
)

const base_uri = "http://www.albinoblacksheep.com/audio/midi/"

func get_download_uri(slug string) string {
	midi_resp, _ := http.Get(base_uri + slug)
	slug_tokens := html.NewTokenizer(midi_resp.Body)
	for {
		if slug_tokens.Next() == html.ErrorToken {
			// Returning io.EOF indicates success.
			return ""
		}
		tok := slug_tokens.Token()
		if tok.Data == "embed" {
			file_url := webutils.GetAttr("src", tok.Attr)
			if len(file_url) > 0 {
				return base_uri + file_url
			}
		}
	}
	return ""
}

func GetSongsBody() io.ReadCloser {
	resp, _ := http.Get(base_uri)
	return resp.Body
}

func DownloadMids(urls []string, songs_dir string) {
	for i := 0; i < len(urls); i++ {
		url := urls[i]
		tokens := strings.Split(url, "/")
		fileName := tokens[len(tokens)-1]
		webutils.DownloadFromUrl(url, songs_dir, fileName)
	}
}

func ScrapeMids(songs_body io.ReadCloser) []string {
	var urls []string

	z := html.NewTokenizer(songs_body)

	for {
		if z.Next() == html.ErrorToken {
			// Returning io.EOF indicates success.
			break
		}
		tok := z.Token()
		if tok.Data == "option" {
			for i := 0; i < len(tok.Attr); i++ {
				download_slug := webutils.GetAttr("value", tok.Attr)
				download_uri := get_download_uri(download_slug)
				urls = append(urls, download_uri)
			}
		}
	}
	return urls
}
