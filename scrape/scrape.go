package scrape

import (
	"fmt"
	"golang.org/x/net/html"
	"io"
	"os"
	"strings"
	// "io/ioutil"
	"net/http"
	"path"
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
			file_url := get_attr("src", tok.Attr)
			if len(file_url) > 0 {
				return base_uri + file_url
			}
		}
	}
	return ""
}

func get_attr(key string, attrs []html.Attribute) string {
	for i := 0; i < len(attrs); i++ {
		if attrs[i].Key == key {
			return attrs[i].Val
		}
	}
	return ""
}

func downloadFromUrl(url, dest_dir string) {
	// Took this shamelessly from : https://github.com/thbar/golang-playground/blob/master/download-files.go
	tokens := strings.Split(url, "/")
	fileName := path.Join(dest_dir, tokens[len(tokens)-1])
	fmt.Println("Downloading", url, "to", fileName)

	// TODO: check file existence first with io.IsExist
	output, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error while creating", fileName, "-", err)
		return
	}
	defer output.Close()

	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return
	}
	defer response.Body.Close()

	n, err := io.Copy(output, response.Body)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return
	}

	fmt.Println(n, "bytes downloaded.")
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
				download_slug := get_attr("value", tok.Attr)
				download_uri := get_download_uri(download_slug)
				downloadFromUrl(download_uri, songs_dir)
			}
		}
	}
}
