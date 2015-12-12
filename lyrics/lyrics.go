package lyrics

import (
	"bufio"
	"fmt"
	"github.com/tebriel/cli-goke/webutils"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"io"
	"log"
	"os"
	"path"
	"regexp"
	"strings"
	"time"
)

const BaseUri = "http://search.azlyrics.com/search.php?q="
const TOSToken = "<!-- Usage of azlyrics.com content by any " +
	"third-party lyrics provider is prohibited by " +
	"our licensing agreement. Sorry about that. -->"

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
	log.Println("Scraping for top link")
	z := html.NewTokenizer(search_body)
	td_started := false
	for {
		if z.Next() == html.ErrorToken {
			// Returning io.EOF indicates success.
			break
		}
		tok := z.Token()
		if tok.DataAtom == atom.Td && tok.Type == html.StartTagToken {
			td_started = true
		} else if tok.DataAtom == atom.A && tok.Type == html.StartTagToken {
			if td_started {
				href, _ := webutils.GetAttr("href", tok.Attr)
				return href
			}
		}
	}
	return ""
}

func ScrapeLyricsFromPage(lyrics_body io.ReadCloser) []string {
	log.Println("Scraping Lyrics from Page")
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

		if in_lyrics && tok.Type == html.EndTagToken && tok.DataAtom == atom.Div {
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
	log.Printf("Getting top lyrics result link for %s\n", query)
	search_body := webutils.GetWebBody(song_query)
	return ScrapeTopLink(search_body)
}

func LyricsAlreadyDownloaded(song_file, lyrics_dir string) bool {
	save_name := CreateSaveName(song_file)
	return webutils.FileExists(path.Join(lyrics_dir, save_name))
}

func ScrapeLyrics(song_file, lyrics_dir string) {
	save_name := CreateSaveName(song_file)
	if !LyricsAlreadyDownloaded(song_file, lyrics_dir) {
		top_link := GetTopLink(song_file)
		if top_link == "" {
			return
		}
		lyrics_page := webutils.GetWebBody(top_link)
		lyrics := ScrapeLyricsFromPage(lyrics_page)
		lyrics_file, _ := os.Create(path.Join(lyrics_dir, save_name))
		for _, lyric := range lyrics {
			lyrics_file.WriteString(lyric)
		}
	} else {
		file, err := os.Open(path.Join(lyrics_dir, save_name))
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
			time.Sleep(1 * time.Second)
		}
	}
}
