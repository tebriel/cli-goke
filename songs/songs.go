package songs

import (
	"fmt"
	"github.com/tebriel/cli-goke/webutils"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"io"
	"io/ioutil"
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
			file_url, err := webutils.GetAttr("src", tok.Attr)
			if err == nil && file_url != "" {
				return file_url
			}
		}
	}
	return ""
}

func DownloadMids(urls []string, songs_dir string) {
	for _, url := range urls {
		tokens := strings.Split(url, "/")
		fileName := tokens[len(tokens)-1]
		webutils.DownloadFromUrl(url, songs_dir, fileName)
	}
}

func DoItAll(songs_dir string) {
	var urls []string
	songs_body := webutils.GetWebBody(BaseUri)
	slugs := ScrapeSlugs(songs_body)
	for _, slug := range slugs {
		mid_body := webutils.GetWebBody(BaseUri + slug)
		urls = append(urls, ScrapeMidUrl(mid_body))
	}
	DownloadMids(urls, songs_dir)
}

func PrintAllSongs(songs_dir string) {
	file_infos, _ := ioutil.ReadDir(songs_dir)
	for _, a_file := range file_infos {
		if a_file.IsDir() {
			continue
		}

		fmt.Println(a_file.Name())
	}
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
			slug, err := webutils.GetAttr("value", tok.Attr)
			// Happens when there was no attr found
			if err == nil && slug != "" {
				slugs = append(slugs, slug)
			}
		}
	}
	return slugs
}
