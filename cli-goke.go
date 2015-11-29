package main

import (
	"github.com/tebriel/cli-goke/lyrics"
	"github.com/tebriel/cli-goke/songs"
	"github.com/tebriel/cli-goke/webutils"
	"os"
	"os/user"
	"path"
)

func main() {
	cur_user, _ := user.Current()
	homedir := cur_user.HomeDir
	songs_dir := path.Join(homedir, ".cliaoke", "songs")
	lyrics_dir := path.Join(homedir, ".cliaoke", "lyrics")
	if !webutils.FileExists(songs_dir) {
		os.MkdirAll(songs_dir, os.ModeDir|0755)
		songs.DoItAll(songs_dir)
	}

	if !webutils.FileExists(lyrics_dir) {
		os.MkdirAll(lyrics_dir, os.ModeDir|0755)
	}
	lyrics.ScrapeLyrics("2Pac_-_California.mid", lyrics_dir)
}
