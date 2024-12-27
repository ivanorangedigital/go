package services

import (
	"digitalcorporation/pkg/utils"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
)

type image struct {
	Name  string `json:"name"`
	Src   string `json:"src"`
	Error string `json:"error"`
}

type ImageUploader struct {
	RootDir           string
	PrefixRequest     string
	MaxLength         int
	MaxSizeFile       int
	AllowedExtensions []string
}

func (u *ImageUploader) Upload(r *http.Request, fieldForm string) ([]*image, error) {
	if err := r.ParseMultipartForm(int64(u.MaxSizeFile) * int64(u.MaxLength)); err != nil {
		return nil, err
	}

	// get files from multipartForm
	files, ok := r.MultipartForm.File[fieldForm]
	if !ok {
		return nil, fmt.Errorf("No file/s provided to field: %s", fieldForm)
	}
	if len(files) > u.MaxLength {
		return nil, fmt.Errorf("You can upload up to a maximum of %d files, current: %d", u.MaxLength, len(files))
	}

	// generate folder struct from date
	generatedStruct := utils.GenerateFolderStructFromDate()
	dir := path.Join(u.RootDir, generatedStruct)

	// create dir if does not exist
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return nil, err
	}

	// result handler
	res := []*image{}

	// iterate files
	for _, fileHeader := range files {
		// initialize image struct
		image := new(image)

		// set file name
		image.Name = fileHeader.Filename

		// check size
		if fileHeader.Size > int64(u.MaxSizeFile) {
			image.Error = fmt.Sprintf("This file exceeds the maximum upload size of %dmb", u.MaxSizeFile)
			res = append(res, image)
			continue
		}

		// get extension
		splittedName := strings.Split(image.Name, ".")
		extension := splittedName[len(splittedName)-1]

		// --- ext (check if ext is allowed) ---
		isExtensionAllowed := false
		for _, ext := range u.AllowedExtensions {
			if extension == ext {
				isExtensionAllowed = true
				break
			}
		}
		if !isExtensionAllowed {
			image.Error = fmt.Sprintf("This extension is not supported, below are the supported extensions: %v", u.AllowedExtensions)
			res = append(res, image)
			continue
		}
		// --- ext ---

		// --- fileName (generate random fileName) ---
		fileName := fmt.Sprintf("%s.%s", utils.GenerateRandomUint(15), extension)

		// --- filePath (retrieve full path, [relative to current execution]) ---
		filePath := path.Join(dir, fileName)

		// --- openFile (open passed file) ---
		file, err := fileHeader.Open()
		if err != nil {
			image.Error = err.Error()
			res = append(res, image)
			continue
		}
		defer file.Close()

		newFile, err := os.Create(filePath)
		if err != nil {
			image.Error = err.Error()
			res = append(res, image)
			continue
		}
		defer newFile.Close()

		if _, err = io.Copy(newFile, file); err != nil {
			image.Error = err.Error()
			res = append(res, image)
			continue
		}

		image.Src = fmt.Sprintf("%s%s%s", u.PrefixRequest, generatedStruct, fileName)
		res = append(res, image)
	}

	return res, nil
}
