package songs

import (
	"os"
	"testing"
)

func TestScrapeMids(t *testing.T) {
	file, err := os.Open("./testfiles/songs.html")
	if err != nil {
		t.Error("Couldn't open test html file")
	}
	defer file.Close()
	urls := ScrapeMids(file)
	if len(urls) != 151 {
		t.Error("There were not 151 urls extracted from the test file, found: ", len(urls))
	}
}
