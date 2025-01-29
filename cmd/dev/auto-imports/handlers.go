// this file watch [x] dir and from them auto import the defined handlers in main
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
	startFlag       = "// -- handlers"
	endFlag         = "// -- end handlers"
	applicationRoot = "cmd/web/"
	entryFile       = applicationRoot + "main.go"
	dirToWatch      = applicationRoot + "handlers"
	projectName     = "digitalcorporation"
)

func base(name string) string {
	return fmt.Sprintf("_ \"%s/%s\"", path.Join(projectName, dirToWatch), name)
}

func main() {
	// open file
	file, err := os.Open(entryFile)
	if err != nil {
		log.Fatal(err)
	}

	// read file
	bts, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
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

	// perform copy
	newFile := make([]byte, len(bts[:indexStart+len(startFlag)+len("\n")]))
	copy(newFile, bts[:indexStart+len(startFlag)+len("\n")])

	for _, dir := range dirs {
		if !dir.IsDir() {
			continue
		}

		fmt.Println("package", dir.Name(), "(handler) imported")

		// append import
		newFile = append(newFile, []byte("\t"+base(dir.Name())+"\n")...)
	}

	// append rest of file
	newFile = append(newFile, bts[indexEnd-len("\t"):]...)

	// close file and reopen with truncate and rdw method
	file.Close()
	file, err = os.OpenFile(entryFile, os.O_RDWR|os.O_TRUNC, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	_, err = io.Copy(file, bytes.NewReader(newFile))
	if err != nil {
		log.Fatal(err)
	}
}
