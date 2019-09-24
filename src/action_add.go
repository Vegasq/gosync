package main

import (
	"fmt"
	"os"
)

func addFileCli(ctx *context){
	for _, sf := range ctx.fileList {
		if sf.Source == ctx.requestedParam {
			fmt.Printf("%s File already exists.\n", emojiSadFace)
			os.Exit(1)
		}
	}

	fmt.Printf("%s Adding file %s\n", emojiSandClock, ctx.requestedParam)

	path := createFileBackupPath(ctx.requestedParam)
	storedFile := copyToStorage(ctx.requestedParam, path)

	ctx.fileList = append(ctx.fileList, storedFile)
}