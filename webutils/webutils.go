package webutils

import (
	"fmt"
	"golang.org/x/net/html"
	"io"
	"net/http"
	"os"
	"path"
)

func GetAttr(key string, attrs []html.Attribute) string {
	for i := 0; i < len(attrs); i++ {
		if attrs[i].Key == key {
			return attrs[i].Val
		}
	}
	return ""
}

func GetWebBody(uri string) io.ReadCloser {
	resp, _ := http.Get(uri)
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
