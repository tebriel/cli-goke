package lyrics

import (
	"os"
	"testing"
)

const song_title = "Will_Smith_-_Wild_Wild_West.mid"

func TestCleanSongName(t *testing.T) {
	actual := CleanSongName(song_title)
	expected := "Will+Smith+++Wild+Wild+West"

	if actual != expected {
		t.Errorf("Expected: %s, Actual: %s", expected, actual)
	}
}

func TestCreateSaveName(t *testing.T) {
	actual := CreateSaveName(song_title)
	expected := "Will_Smith_-_Wild_Wild_West.txt"

	if actual != expected {
		t.Errorf("Expected: %s, Actual: %s", expected, actual)
	}
}

func TestScrapeTopLink(t *testing.T) {
	file, err := os.Open("./testfiles/wil_smith_wild_wild_west.html")
	if err != nil {
		t.Error("Couldn't open test html file")
	}
	defer file.Close()
	actual := ScrapeTopLink(file)
	expected := "http://www.azlyrics.com/lyrics/willsmith/wildwildwest.html"
	if actual != expected {
		t.Errorf("Scraped the wrong url, expected: %s, actual: %s", expected, actual)
	}
}
