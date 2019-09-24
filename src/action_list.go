package main

import "fmt"

func displayFile(sf storedFile){
	fmt.Printf("%s %s\t%s\t%s %s\n", emojiFile, sf.Source, sf.FileMode, emojiClock, sf.Updated)
}
func displayFiles(sfs storedFiles){
	fmt.Printf("%s List of files\n", emojiBookStack)

	for i := range sfs {
		displayFile(sfs[i])
	}
}

func listCli(ctx *context) {
	displayFiles(ctx.fileList)
}