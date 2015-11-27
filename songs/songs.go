package songs

import (
	"github.com/tebriel/cli-goke/webutils"
	"golang.org/x/net/html"
	// "io/ioutil"
	"net/http"
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

func ScrapeMids(songs_dir string) {
	resp, _ := http.Get(base_uri)
	// body, _ := ioutil.ReadAll(resp.Body)
	// fmt.Println(string(body))
	z := html.NewTokenizer(resp.Body)
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
				webutils.DownloadFromUrl(download_uri, songs_dir)
			}
		}
	}
}
