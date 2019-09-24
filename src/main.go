package main

import (
	"log"
	"os"
	"strings"
)

const ADD = "add"
const RM = "rm"
const PUSH = "push"
const PULL = "pull"
const LIST = "list"

type context struct {
	requestedAction string
	requestedParam string
	fileList storedFiles

	git GitStorage
}


func createFileBackupPath(backupFilePath string) string {
	pathNoFile := strings.Split(backupFilePath, "/")
	pathNoFile = pathNoFile[0:len(pathNoFile)-1]

	path := strings.Join(pathNoFile, "/")

	backupPath := storageDir + path
	err := os.MkdirAll(backupPath, os.ModeDir | os.ModePerm)
	if err != nil {
		log.Fatalln("Failed to create dir", backupPath, err)
	}

	return storageDir + backupFilePath
}

func saveConfig(ctx *context){
	saveStorage(ctx.fileList)
}

func isValidCommand(cmd string) bool {
	availableCommands := []string{ADD, RM, PULL, PUSH, LIST}

	for i := range availableCommands {
		if cmd == availableCommands[i] {
			return true
		}
	}

	return false
}


func executeAction(ctx *context) {
	if ctx.requestedAction == ADD {
		addFileCli(ctx)
		saveConfig(ctx)
	} else if ctx.requestedAction == RM {
		rmFileCli(ctx)
		saveConfig(ctx)
	} else if ctx.requestedAction == PULL {
		pullCli(ctx)
	} else if ctx.requestedAction == PUSH {
		pushCli(ctx)
	} else if ctx.requestedAction == LIST {
		listCli(ctx)
	}
}

func main() {
	context := initApp()
	parseArguments(context)

	executeAction(context)
}
