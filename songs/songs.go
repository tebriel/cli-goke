package songs

import (
	"github.com/tebriel/cli-goke/webutils"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"io"
	"strings"
)

const BaseUri = "http://www.albinoblacksheep.com/audio/midi/"

func ScrapeMidUrl(mid_body io.ReadCloser) string {
	slug_tokens := html.NewTokenizer(mid_body)
	for {
		if slug_tokens.Next() == html.ErrorToken {
			// Returning io.EOF indicates success.
			return ""
		}
		tok := slug_tokens.Token()
		if tok.DataAtom == atom.Embed {
			file_url := webutils.GetAttr("src", tok.Attr)
			if len(file_url) > 0 {
				return file_url
			}
		}
	}
	return ""
}

func DownloadMids(urls []string, songs_dir string) {
	for i := 0; i < len(urls); i++ {
		url := urls[i]
		tokens := strings.Split(url, "/")
		fileName := tokens[len(tokens)-1]
		webutils.DownloadFromUrl(url, songs_dir, fileName)
	}
}

func DoItAll(songs_dir string) {
	var urls []string
	songs_body := webutils.GetWebBody(BaseUri)
	slugs := ScrapeSlugs(songs_body)
	for i := 0; i < len(slugs); i++ {
		mid_body := webutils.GetWebBody(BaseUri + slugs[i])
		urls = append(urls, ScrapeMidUrl(mid_body))
	}
	DownloadMids(urls, songs_dir)
}

func ScrapeSlugs(songs_body io.ReadCloser) []string {
	var slugs []string

	z := html.NewTokenizer(songs_body)

	for {
		if z.Next() == html.ErrorToken {
			// Returning io.EOF indicates success.
			break
		}
		tok := z.Token()
		if tok.DataAtom == atom.Option {
			for i := 0; i < len(tok.Attr); i++ {
				slugs = append(slugs, webutils.GetAttr("value", tok.Attr))
				// download_slug := webutils.GetAttr("value", tok.Attr)
				// download_uri := get_download_uri(download_slug)
				// urls = append(urls, download_uri)
			}
		}
	}
	return slugs
}
