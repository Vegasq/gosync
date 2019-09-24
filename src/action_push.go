package main

import "fmt"

func pushCli(ctx *context){
	for i := range ctx.fileList {
		if len(ctx.fileList[i].Source) == 0 {
			continue
		}
		fmt.Printf("%s Pushing file %s\n", emojiSave, ctx.fileList[i].Source)
		path := createFileBackupPath(ctx.fileList[i].Storage)
		//copy(ctx.fileList[i].Source, path)
		copyToStorage(ctx.fileList[i].Source, path)
	}

	ctx.git.commit()
	ctx.git.push()
	fmt.Printf("%s Done\n", emojiCool)
}
