package main

import (
	"github.com/tebriel/cli-goke/scrape"
	"os"
	"os/user"
	"path"
)

// Taken from http://stackoverflow.com/questions/10510691/how-to-check-whether-a-file-or-directory-denoted-by-a-path-exists-in-golang
// exists returns whether the given file or directory exists or not
func exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func main() {
	cur_user, _ := user.Current()
	homedir := cur_user.HomeDir
	songs_dir := path.Join(homedir, ".cliaoke", "songs")
	lyrics_dir := path.Join(homedir, ".cliaoke", "lyrics")
	if !exists(songs_dir) {
		os.MkdirAll(songs_dir, os.ModeDir|0755)
		scrape.ScrapeMids(songs_dir)
	}

	if !exists(lyrics_dir) {
		os.MkdirAll(lyrics_dir, os.ModeDir|0755)
	}
}
