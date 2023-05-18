package controller

import (
	"fmt"
	"go-mock/model"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/mbiagini/go-server-utils/gsrender"
)

func HandleRequest(w http.ResponseWriter, r *http.Request, e model.Endpoint) {

	resp, err := e.FindResponse(r)
	if err != nil {
		gsrender.WriteJSON(w, http.StatusInternalServerError, model.ErrorFrom(err))
		return
	}

	if resp.Delay > 0 {
		time.Sleep(time.Duration(resp.Delay) * time.Millisecond)
	}

	if resp.HasBody() {
		w.Header().Set("Content-type", *resp.ContentType)
		w.WriteHeader(resp.Code)
		buf, err := readFile(*resp.BodyFilename)
		if err != nil {
			gsrender.WriteJSON(w, http.StatusInternalServerError, model.ErrorFrom(err))
			return
		}
		w.Write(buf)
	} else {
		w.WriteHeader(resp.Code)
	}

}

func readFile(f string) ([]byte, error) {
	buf, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, fmt.Errorf("io error reading file %s: %s", f, err.Error())
	}
	return buf, nil
}