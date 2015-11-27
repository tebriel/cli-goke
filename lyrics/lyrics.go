package lyrics

import (
	"github.com/tebriel/cli-goke/webutils"
	"golang.org/x/net/html"
	"strings"
	// "io/ioutil"
	"net/http"
	"regexp"
)

const base_uri = "ttp://search.azlyrics.com/search.php?q="

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

func clean_song_name(song_file string) string {
	non_chars := regexp.MustCompile("[^a-zA-Z0-9\n.]")
	// def clean_input(song_file):
	// song = song_file.replace('.mid', '').replace('.kar', '')
	// song = re.sub('[^a-zA-Z0-9\n\.]', '+', song)
	// return song
	result := strings.Replace(song_file, ".mid", "", -1)
	result = strings.Replace(result, ".kar", "", -1)
	result = non_chars.ReplaceAllString(result, "+")
	return result
}

func get_top_link(song_name string) string {
	query := clean_song_name(song_name)
	song_query := base_uri + query
	resp, _ := http.Get(song_query)
	z := html.NewTokenizer(resp.Body)
	for {
		if z.Next() == html.ErrorToken {
			// Returning io.EOF indicates success.
			break
		}
		tok := z.Token()
		if tok.Data == "td" {
			// Do some logic here to check for start/end token and search for
			// the lyrics inside of it
		}
	}
	return "alink"
}

func ScrapeLyrics(song_file, lyrics_dir string) {
	lyrics_file := strings.Replace(song_file, ".mid", ".txt", -1)
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
				webutils.DownloadFromUrl(download_uri, lyrics_dir, lyrics_file)
			}
		}
	}
}
