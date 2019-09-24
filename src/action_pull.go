package main

func pullCli(ctx *context){
	for i := range ctx.fileList {
		if len(ctx.fileList[i].Source) == 0 {
			continue
		}

		copyFromStorage(ctx.fileList[i])
	}

	ctx.git.pull()
}
