package main

import (
	"github.com/tebriel/cli-goke/scrape"
	"os"
	"os/user"
	"path"
)

func main() {
	cur_user, _ := user.Current()
	homedir := cur_user.HomeDir
	songs_dir := path.Join(homedir, ".cliaoke", "songs")
	os.MkdirAll(songs_dir, os.ModeDir|0755)
	scrape.ScrapeMids(songs_dir)
}
