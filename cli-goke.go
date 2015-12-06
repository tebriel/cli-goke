package main

import (
	"flag"
	"github.com/tebriel/cli-goke/lyrics"
	"github.com/tebriel/cli-goke/songs"
	"github.com/tebriel/cli-goke/webutils"
	"os"
	"os/user"
	"path"
)

var list_songs bool
var sing_song string

func init_flags() {
	flag.BoolVar(&list_songs, "songs", false, "Pass this flag to list all the songs")
	flag.StringVar(&sing_song, "sing", "", "Pass the name of a song here to sing it")
	flag.Parse()

	if !list_songs && sing_song == "" {
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

	if !webutils.FileExists(lyrics_dir) {
		os.MkdirAll(lyrics_dir, os.ModeDir|0755)
	}

	if list_songs {
		songs.PrintAllSongs(songs_dir)
		os.Exit(0)
	} else {
		lyrics.ScrapeLyrics(sing_song, lyrics_dir)
	}
}
