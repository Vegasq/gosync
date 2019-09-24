package main

import (
	"fmt"
	"log"
	"os/user"
)
import "gopkg.in/gcfg.v1"


var CONFIG = struct {
	Main struct {
		Repo string
	}
}{}

var storageDir string
var dbPath string
var configPath string


const emojiFile = "📃"
const emojiClock = "🕐"
const emojiBookStack = "📚"
const emojiSadFace = "😢"
const emojiSandClock = "⏳"
const emojiHammer = "🔨"
const emojiCool = "👍"
const emojiSave = "💾"

func parseConfig() {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	storageDir = usr.HomeDir + "/.config/gosync/store/"
	dbPath = usr.HomeDir + "/.config/gosync/db"
	configPath = usr.HomeDir + "/.config/gosync/config"

	err = gcfg.ReadFileInto(&CONFIG, configPath)

	if err != nil {
		fmt.Printf("Config error: %s\n", err)
	}
}
