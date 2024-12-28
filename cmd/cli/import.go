package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

var startFlag = "// -- handlers"
var endFlag = "// -- end handlers"

func base(name string) string {
	return fmt.Sprintf("_ \"digitalcorporation/cmd/web2/handlers/%s\"", name)
}

func main() {
	file, err := os.Open("./cmd/web2/main.go")
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
		log.Fatal(errors.New("Unable to find '// -- handlers' flag"))
	}

	// detect line '// -- end handlers'
	indexEnd := bytes.Index(bts, []byte(endFlag))
	if indexEnd == -1 {
		log.Fatal(errors.New("Unable to find '// -- end handlers' flag"))
	}

	dirs, err := os.ReadDir("./cmd/web2/handlers")
	if err != nil {
		log.Fatal(err)
	}

	newFile := []byte{}
	newFile = append(newFile, bts[:indexStart+len(startFlag)]...)

	for _, dir := range dirs {
		if !dir.IsDir() {
			continue
		}

		// append import
		newFile = append(newFile, []byte("\n\t"+base(dir.Name())+"\n")...)
	}

	// append rest of file
	newFile = append(newFile, bts[indexEnd-len("\t"):]...)

	// close file and reopen with truncate and rdw method
	file.Close()
	file, err = os.OpenFile("./cmd/web2/main.go", os.O_RDWR|os.O_TRUNC, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	_, err = io.Copy(file, bytes.NewReader(newFile))
	if err != nil {
		log.Fatal(err)
	}
}
