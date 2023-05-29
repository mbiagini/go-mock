package controller

import (
	"fmt"
	"go-mock/apierrors"
	"go-mock/service"
	"net/http"

	"github.com/mbiagini/go-server-utils/gslog"

	"github.com/mbiagini/go-server-utils/gsrender"
)

func GetAllFiles(w http.ResponseWriter, r *http.Request) {

	response, err := service.FindAllFiles()
	if err != nil {
		gsrender.WriteJSON(w, http.StatusInternalServerError, apierrors.New(apierrors.INTERNAL_SERVER_ERROR))
		return
	}

	gsrender.WriteJSON(w, http.StatusOK, response)
}

func PostFiles(w http.ResponseWriter, r *http.Request, uploadMaxSize int) {

	// Set maximum multipart form size.
	maxFormBytes := int64(uploadMaxSize) << 20
	err := r.ParseMultipartForm(maxFormBytes)
	if err != nil {
		gsrender.WriteJSON(w, http.StatusRequestEntityTooLarge, apierrors.New(apierrors.FILES_SIZE_EXCEEDED))
		return
	}

	files := r.MultipartForm.File["files"]
	for _, header := range files {
		file, err := header.Open()
		if err != nil {
			gsrender.WriteJSON(w, http.StatusInternalServerError, apierrors.New(apierrors.IO_FILE_ERROR))
			return
		}
		defer file.Close()

		gslog.Server(fmt.Sprintf("Received file with name: %s (%d bytes)", header.Filename, header.Size))

		err = service.UploadFile(header.Filename, file)
		if err != nil {
			gsrender.WriteJSON(w, http.StatusInternalServerError, apierrors.New(apierrors.ERR_SAVING_FILE))
			return
		}
	}

	gsrender.Status(w, http.StatusCreated)
}