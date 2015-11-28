package lyrics

import (
	"github.com/tebriel/cli-goke/webutils"
	"golang.org/x/net/html"
	"io"
	"net/http"
	"regexp"
	"strings"
)

const BaseUri = "http://search.azlyrics.com/search.php?q="
const TOSToken = "<!-- Usage of azlyrics.com content by any " +
	"third-party lyrics provider is prohibited by " +
	"our licensing agreement. Sorry about that. -->"

func get_download_uri(slug string) string {
	midi_resp, _ := http.Get(BaseUri + slug)
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
				return BaseUri + file_url
			}
		}
	}
	return ""
}

func CleanSongName(song_file string) string {
	non_chars := regexp.MustCompile("[^a-zA-Z0-9\n.]")
	result := strings.Replace(song_file, ".mid", "", -1)
	result = strings.Replace(result, ".kar", "", -1)
	result = non_chars.ReplaceAllString(result, "+")
	return result
}

func CreateSaveName(song_file string) string {
	return strings.Replace(song_file, ".mid", ".txt", -1)
}

func ScrapeTopLink(search_body io.ReadCloser) string {
	z := html.NewTokenizer(search_body)
	td_started := false
	for {
		if z.Next() == html.ErrorToken {
			// Returning io.EOF indicates success.
			break
		}
		tok := z.Token()
		if tok.Data == "td" && tok.Type == html.StartTagToken {
			td_started = true
		} else if tok.Data == "a" && tok.Type == html.StartTagToken {
			if td_started {
				href := webutils.GetAttr("href", tok.Attr)
				return href
			}
		}
	}
	return ""
}

func ScrapeLyricsFromPage(lyrics_body io.ReadCloser) []string {
	var result []string
	z := html.NewTokenizer(lyrics_body)
	in_lyrics := false
	for {
		if z.Next() == html.ErrorToken {
			// Returning io.EOF indicates success.
			break
		}
		tok := z.Token()
		if tok.Type == html.CommentToken {
			if string(z.Raw()) == TOSToken {
				in_lyrics = true
				continue
				// Yay we found the lyrics
			}
		}

		if in_lyrics && tok.Type == html.EndTagToken && tok.Data == "div" {
			return result
		}

		if in_lyrics {
			if tok.Type == html.TextToken {
				result = append(result, string(z.Raw()))
			}
		}
	}
	return result
}

func GetTopLink(song_name string) string {
	query := CleanSongName(song_name)
	song_query := BaseUri + query
	webutils.GetWebBody(song_query)
	return "alink"
}

func ScrapeLyrics(song_file, lyrics_dir string) {
	lyrics_file := CreateSaveName(song_file)
	resp, _ := http.Get(BaseUri)
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
				webutils.DownloadFromUrl(download_uri, lyrics_dir, lyrics_file)
			}
		}
	}
}
