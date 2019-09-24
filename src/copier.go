package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"
)

type storedFile struct {
	Source   string
	Storage  string
	FileMode os.FileMode

	//Added   time.Time
	Updated time.Time
}

func fileToFileMode(fl *os.File) os.FileMode {
	fi, err := fl.Stat()
	if err != nil {
		log.Fatalf("Failed to read stat from file %s with %s.\n", fl, err)
	}

	return fi.Mode()
}

func openFiles(from string, to string) (*os.File, *os.File) {
	// Try to open Source file
	flFrom, err := os.Open(from)
	if err != nil {
		log.Fatalf("File %s to backup not found. %s\n", from, err)
	}

	// Try to create destination file
	fileTo, err := os.Create(to)
	if err != nil {
		log.Fatalf("Failed to create %s with %s\n", to, err)
	}

	return flFrom, fileTo
}

func copyToStorage(from string, to string) storedFile {
	fromFile, toFile := openFiles(from, to)
	defer fromFile.Close()
	defer toFile.Close()

	binData, readAllErr := ioutil.ReadAll(fromFile)
	if readAllErr != nil {
		log.Printf("Failed to read from file %s with %s.\n", from, readAllErr)
	}

	_, writeFileErr := toFile.Write(binData)
	if writeFileErr != nil {
		fmt.Printf("Failed to write to file %s with %s\n", to, writeFileErr)
	}

	return storedFile {
		Source:   from,
		Storage:  to,
		FileMode: fileToFileMode(fromFile),
		Updated:  time.Now(),
	}
}

func copyFromStorage(sf storedFile) {
	source, storage := openFiles(sf.Source, sf.Storage)
	defer source.Close()
	defer storage.Close()

	binData, readAllErr := ioutil.ReadAll(storage)
	if readAllErr != nil {
		log.Printf("Failed to read from file %s with %s.\n", sf.Storage, readAllErr)
	}

	_, writeFileErr := source.WriteAt(binData, 0)
	if writeFileErr != nil {
		fmt.Printf("Failed to write to file %s with %s\n", source.Name(), writeFileErr)
	}

	chmodErr := source.Chmod(sf.FileMode)
	if chmodErr != nil {
		log.Fatalln("Failed to set chmod on file", source.Name())
	}
}
