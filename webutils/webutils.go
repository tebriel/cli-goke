package webutils

import (
	"fmt"
	"golang.org/x/net/html"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
)

func GetAttr(key string, attrs []html.Attribute) string {
	for i := 0; i < len(attrs); i++ {
		if attrs[i].Key == key {
			return attrs[i].Val
		}
	}
	return ""
}

func DownloadFromUrl(url, dest_dir string) {
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
