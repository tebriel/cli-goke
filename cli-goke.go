package main

import (
	// "github.com/tebriel/cli-goke/lyrics"
	"flag"
	"github.com/tebriel/cli-goke/songs"
	"github.com/tebriel/cli-goke/webutils"
	"os"
	"os/user"
	"path"
)

var songsvar bool
var sing_song string

func init_flags() {
	flag.BoolVar(&songsvar, "songs", false, "Pass this flag to list all the songs")
	flag.StringVar(&sing_song, "sing", "", "Pass the name of a song here to sing it")
	flag.Parse()

	if !songsvar && sing_song == "" {
		flag.PrintDefaults()
	}
}

func main() {
	init_flags()

	cur_user, _ := user.Current()
	homedir := cur_user.HomeDir
	songs_dir := path.Join(homedir, ".cliaoke", "songs")
	lyrics_dir := path.Join(homedir, ".cliaoke", "lyrics")

	if !webutils.FileExists(songs_dir) {
		os.MkdirAll(songs_dir, os.ModeDir|0755)
		songs.DoItAll(songs_dir)
	}

	if songsvar {
		songs.PrintAllSongs(songs_dir)
		os.Exit(0)
	}

	if !webutils.FileExists(lyrics_dir) {
		os.MkdirAll(lyrics_dir, os.ModeDir|0755)
		// lyrics.ScrapeLyrics("2Pac_-_California.mid", lyrics_dir)
	}
}
