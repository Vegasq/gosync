package main

import (
	"fmt"
	"os"
)

func createStorageDir(){
	mkdirErr := os.Mkdir(storageDir, os.ModeDir | os.ModePerm)
	if os.IsExist(mkdirErr) == false {
		fmt.Printf("%s Failed to mkdir %s \n", emojiSadFace, mkdirErr)
	}
}

func initStorage() {
	createStorageDir()

	isNew := initStorageCreate()
	if isNew {
		saveStorage(storedFiles{})
	}
}

func initStorageCreate() bool {
	isNew := false

	_, err := os.Open(dbPath)
	if err != nil {
		_, confCreateErr := os.Create(dbPath)
		if confCreateErr != nil {
			fmt.Printf("%s Can't create config", emojiSadFace)
		}
		isNew = true
	}
	return isNew

}

func NewContext() *context {
	ctx := context{}
	ctx.fileList = loadStorage()

	gs := GitStorage{}
	gs.initRepo()

	ctx.git = gs

	return &ctx
}

func initApp () *context {
	parseConfig()
	initStorage()

	// Prepare app
	ctx := NewContext()

	return ctx
}

func parseArguments(ctx *context){
	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Println("Not enough params")
		os.Exit(1)
	}

	if isValidCommand(args[0]) == false {
		fmt.Printf("Unknown command %s\n", args[0])
		os.Exit(1)
	}

	if len(args) == 1 {
		ctx.requestedAction = args[0]
	} else if len(args) == 2 {
		ctx.requestedAction = args[0]
		ctx.requestedParam = args[1]
	} else {
		fmt.Println("Unknown command")
		os.Exit(1)
	}
}

