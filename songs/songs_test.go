package songs

import (
	"os"
	"testing"
)

func TestScrapeSlugs(t *testing.T) {
	file, err := os.Open("./testfiles/songs.html")
	if err != nil {
		t.Error("Couldn't open test html file")
	}
	defer file.Close()
	urls := ScrapeSlugs(file)
	if len(urls) != 151 {
		t.Error("There were not 151 urls extracted from the test file, found: ", len(urls))
	}
}

func TestScrapeMidUrl(t *testing.T) {
	file, err := os.Open("./testfiles/2pac.html")
	if err != nil {
		t.Error("Couldn't open test html file")
	}
	defer file.Close()
	url := ScrapeMidUrl(file)
	expected := "http://media3.albinoblacksheep.com/albino_midi/2Pac_-_Until_The_End_Of_Time.mid"
	if url != expected {
		t.Error("Url wasn't right expected: ", expected, " actual: ", url)
	}
}
