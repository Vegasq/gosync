package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)


const ADD = "add"
const RM = "rm"
const PUSH = "push"
const PULL = "pull"


func initApp () {
	mkdirErr := os.Mkdir(storageDir, os.ModeDir | os.ModePerm)
	if os.IsExist(mkdirErr) == false {
		fmt.Printf("Failed to mkdir %s \n", mkdirErr)
	}
	_, err := os.Open(dbPath)
	if err != nil {
		_, confCreateErr := os.Create(dbPath)
		if confCreateErr != nil {
			fmt.Printf("Can't create config")
		}
	}
}

type context struct {
	requestedAction string
	requestedParam string
	fileList []string

	git GitStorage
}

func saveSyncFilesList(s []string){
	fl, err := os.Create(dbPath)
	if err != nil {
		fmt.Println("Failed to create file")
	}
	defer fl.Close()

	var toString string
	for _, line := range s {
		toString += fmt.Sprintf("%s\n", line)
	}

	// Remove multiple empty lines
	for ;; {
		if strings.Contains(toString, "\n\n") {
			toString = strings.Replace(toString, "\n\n", "\n", -1)
		} else {
			break
		}
	}
	toString = strings.TrimLeft(toString, "\n")

	fl.Write([]byte(toString))
}

func loadSyncFilesList() []string {
	fl, err := os.Open(dbPath)
	if err != nil {
		fmt.Printf("Can't open %s\n", dbPath)
		return []string{}
	}
	defer fl.Close()

	b, err := ioutil.ReadAll(fl)
	lines := strings.Split(string(b), "\n")

	return lines
}

func addFileCli(ctx *context){
	for _, line := range ctx.fileList {
		if line == ctx.requestedParam {
			fmt.Println("File already exists.")
			return
		}
	}
	fmt.Printf("Adding file %s\n", ctx.requestedParam)
	ctx.fileList = append(ctx.fileList, ctx.requestedParam)
}

func rmFileCli(ctx *context){
	var s []string
	for _, line := range ctx.fileList {
		if line != ctx.requestedParam {
			s = append(s, ctx.requestedParam)
		}
	}
	fmt.Printf("Removing file %s\n", ctx.requestedParam)
	ctx.fileList = s
}

func createFileBackupPath(backupFilePath string) string {
	pathNoFile := strings.Split(backupFilePath, "/")
	pathNoFile = pathNoFile[0:len(pathNoFile)-1]
	//fileName := pathNoFile[len(pathNoFile)-1:]
	path := strings.Join(pathNoFile, "/")

	backupPath := storageDir + path
	err := os.MkdirAll(backupPath, os.ModeDir | os.ModePerm)
	if err != nil {
		fmt.Println(err)
	}

	return storageDir + backupFilePath
}

func copy(from string, to string){
	flFrom, err := os.Open(from)
	b, err := ioutil.ReadAll(flFrom)

	if err != nil {
		fmt.Printf("File %s to backup not found. %s\n", from, err)
	} else {
		defer flFrom.Close()
	}

	flTo, err := os.Create(to)
	if err != nil {
		fmt.Println(err)
	}
	_, err = flTo.Write(b)
	if err != nil {
		fmt.Printf("Failed to write to file %s with %s\n", to, err)
	}
}

func pushCli(ctx *context){
	for i := range ctx.fileList {
		if len(ctx.fileList[i]) == 0 {
			continue
		}
		fmt.Printf("Pushing file %s\n", ctx.fileList[i])
		path := createFileBackupPath(ctx.fileList[i])
		copy(ctx.fileList[i], path)
	}

	ctx.git.commit()
	ctx.git.push()
}

func pullCli(ctx *context){
	for i := range ctx.fileList {
		if len(ctx.fileList[i]) == 0 {
			continue
		}

		copy(storageDir + ctx.fileList[i], ctx.fileList[i])
	}
	ctx.git.pull()
}

func readConfig(ctx *context) {
	ctx.fileList = loadSyncFilesList()
}
func saveConfig(ctx *context){
	saveSyncFilesList(ctx.fileList)
}

func isValidCommand(cmd string) bool {
	availableCommands := []string{ADD, RM, PULL, PUSH}

	for i := range availableCommands {
		if cmd == availableCommands[i] {
			return true
		}
	}

	return false
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

func executeAction(ctx *context) {
	if ctx.requestedAction == ADD {
		addFileCli(ctx)
	} else if ctx.requestedAction == RM {
		rmFileCli(ctx)
	} else if ctx.requestedAction == PULL {
		pullCli(ctx)
	} else if ctx.requestedAction == PUSH {
		pushCli(ctx)
	}
}

func main() {
	// Create config and folder
	parseConfig()

	gs := GitStorage{}
	gs.initRepo()

	initApp()

	// Prepare app
	ctx := context{}
	ctx.git = gs
	parseArguments(&ctx)
	readConfig(&ctx)

	// Run app
	executeAction(&ctx)

	// Save results
	saveConfig(&ctx)
}
