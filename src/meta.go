package main

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

// go binary encoder
// https://stackoverflow.com/questions/28020070/golang-serialize-and-deserialize-back
func toBytes(m storedFiles) []byte {
	b := bytes.Buffer{}
	e := gob.NewEncoder(&b)
	err := e.Encode(m)
	if err != nil { fmt.Println(`failed gob Encode`, err) }
	return b.Bytes()
}

// go binary decoder
// https://stackoverflow.com/questions/28020070/golang-serialize-and-deserialize-back
func fromBytes(by []byte) (storedFiles, error) {
	m := storedFiles{}

	b := bytes.Buffer{}
	b.Write(by)
	d := gob.NewDecoder(&b)
	err := d.Decode(&m)
	if err != nil {
		return nil, errors.New(fmt.Sprintln("Failed to gob decode config file.", err))
	}
	return m, nil
}


type storedFiles []storedFile

func loadStorage() storedFiles {
	dbFile, err := os.Open(dbPath)
	if err != nil {
		log.Fatalf("Failed to open DB %s with %s.\n", dbPath, err)
	}

	b, readErr := ioutil.ReadAll(dbFile)
	if readErr != nil {
		log.Fatalf("Failed to read DB %s with %s.\n", dbPath, readErr)
	}

	var sfs storedFiles
	sfs, decodeErr := fromBytes(b)
	if decodeErr != nil {
		log.Fatalf("Failed to decode DB %s with %s.\n", dbPath, decodeErr)
	}

	return sfs
}

func saveStorage(sfs storedFiles) {dbFile, err := os.OpenFile(dbPath, os.O_RDWR, os.ModeAppend)
	//dbFile, err := os.Open(dbPath)
	if err != nil {
		log.Fatalf("Failed to open DB %s with %s.\n", dbPath, err)
	}

	b64sfs := toBytes(sfs)

	_, writeErr := dbFile.Write([]byte(b64sfs))
	if writeErr != nil {
		log.Fatalf("Faile to write file %s with %s.\n", dbPath, writeErr)
	}
}