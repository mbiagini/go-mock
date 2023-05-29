package service

import (
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"os"

	"github.com/mbiagini/go-server-utils/gslog"
)

func FindAllFiles() ([]string, error) {

	dirPath := "./resources/responses"
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		gslog.Server(fmt.Sprintf("Error when trying to read all response files: %s", err.Error()))
		return nil, err
	}

	response := make([]string, 0, len(files))
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		response = append(response, file.Name())
	}

	return response, nil
}

func UploadFile(filename string, f multipart.File) error {

	// Create a new file in the upload directory.
	filePath := "./resources/responses/" + filename
	file, err := os.Create(filePath)
	if err != nil {
		gslog.Server(fmt.Sprintf("Error when trying to create file in path %s: %s", filePath, err.Error()))
		return err
	}
	defer file.Close()

	// Copy the uploaded file to the newly created file.
	_, err = io.Copy(file, f)
	if err != nil {
		gslog.Server(fmt.Sprintf("Error when trying to copy received file to saved one: %s", err.Error()))
		return err
	}

	return nil
}