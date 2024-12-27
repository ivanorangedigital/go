package utils

import (
	"io"
	"os"
	"time"
)

func LastModifiedDateFile(filePath string) (time.Time, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return time.Time{}, err
	}

	// retrieve last date of modified file (in utc)
	return fileInfo.ModTime().UTC().Truncate(time.Second), nil
}

func ReadFile(filePath string) ([]byte, error) {
	// initialize slice of bytes
	bytes := []byte{}

	f, err := os.Open(filePath)
	if err != nil {
		return bytes, err
	}
	defer f.Close()

	// allocate buffer
	buffer := make([]byte, 1024)

	for {
		bytesReaded, err := f.Read(buffer)

		if bytesReaded > 0 {
			bytes = append(bytes, buffer[:bytesReaded]...)
		}

		if err != nil {
			// readed all file
			if err == io.EOF {
				break
			}

			return nil, err
		}
	}

	return bytes, nil
}
