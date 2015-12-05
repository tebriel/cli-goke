package webutils

import (
	"fmt"
	"golang.org/x/net/html"
	"io"
	"log"
	"net/http"
	"os"
	"path"
)

// errorString is a trivial implementation of error.
type errorString struct {
	s string
}

func (e *errorString) Error() string {
	return e.s
}

// New returns an error that formats as the given text.
func New(text string) error {
	return &errorString{text}
}

// Taken from http://stackoverflow.com/questions/10510691/how-to-check-whether-a-file-or-directory-denoted-by-a-path-exists-in-golang
// exists returns whether the given file or directory exists or not
func FileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func GetAttr(key string, attrs []html.Attribute) (val string, err error) {
	for i := 0; i < len(attrs); i++ {
		if attrs[i].Key == key {
			return attrs[i].Val, nil
		}
	}
	return "", New("No attribute with that name found.")
}

func GetWebBody(uri string) io.ReadCloser {
	resp, err := http.Get(uri)
	if err != nil {
		log.Printf("Error scraping lyrics: %s\n", err)
		return nil
	}
	return resp.Body
}

func DownloadFromUrl(url, dest_dir, dest_file string) {
	// Took this shamelessly from : https://github.com/thbar/golang-playground/blob/master/download-files.go
	fileName := path.Join(dest_dir, dest_file)
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
