package main

import "fmt"

func rmFileCli(ctx *context){
	var newStorage storedFiles
	for _, sf := range ctx.fileList {
		if sf.Source != ctx.requestedParam {
			newStorage = append(newStorage, sf)
		}
	}

	if len(ctx.fileList) != len(newStorage) {
		fmt.Printf("%s Removing file %s\n", emojiHammer, ctx.requestedParam)
		ctx.fileList = newStorage
	} else {
		fmt.Printf("%s File %s not found\n", emojiSadFace, ctx.requestedParam)
	}
}
