package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path"
)

var (
	startFlag   = "// -- handlers"
	endFlag     = "// -- end handlers"
	filePath    = "./cmd/web2/main.go"
	dirToWatch  = "./cmd/web2/handlers"
	projectName = "digitalcorporation"
)

func base(name string) string {
	return fmt.Sprintf("_ \"%s/%s\"", path.Join(projectName, dirToWatch), name)
}

func main() {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}

	buffer := make([]byte, 1024)
	bts := []byte{}

	// file reader
	for {
		readedBytes, err := file.Read(buffer)
		if readedBytes > 0 {
			bts = append(bts, buffer[:readedBytes]...)
		}
		if err != nil {
			if err == io.EOF {
				break
			}

			log.Fatal(err)
		}
	}

	// detect line '// -- handlers'
	indexStart := bytes.Index(bts, []byte(startFlag))
	if indexStart == -1 {
		log.Fatal(errors.New(fmt.Sprintf("Unable to find '%s' flag", startFlag)))
	}

	// detect line '// -- end handlers'
	indexEnd := bytes.Index(bts, []byte(endFlag))
	if indexEnd == -1 {
		log.Fatal(errors.New(fmt.Sprintf("Unable to find '%s' flag", endFlag)))
	}

	dirs, err := os.ReadDir(dirToWatch)
	if err != nil {
		log.Fatal(err)
	}

	newFile := []byte{}
	newFile = append(newFile, bts[:indexStart+len(startFlag)+len("\n")]...)

	for _, dir := range dirs {
		if !dir.IsDir() {
			continue
		}

		// append import
		newFile = append(newFile, []byte("\t"+base(dir.Name())+"\n")...)
	}

	// append rest of file
	newFile = append(newFile, bts[indexEnd-len("\t"):]...)

	// close file and reopen with truncate and rdw method
	file.Close()
	file, err = os.OpenFile(filePath, os.O_RDWR|os.O_TRUNC, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	_, err = io.Copy(file, bytes.NewReader(newFile))
	if err != nil {
		log.Fatal(err)
	}
}
