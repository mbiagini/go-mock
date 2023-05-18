package model

import (
	"errors"
	"fmt"
	"os"
)

type Response struct {
	Id           int     `json:"id"`
	Code         int     `json:"code"`
	ContentType  *string `json:"content_type"`
	BodyFilename *string `json:"body_filename"`
	Delay        int     `json:"delay"`
}

func (r *Response) Validate() error {
	if r.BodyFilename != nil {
		if _, err := os.Stat(*r.BodyFilename); errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("filename %s not found for response with id %d", *r.BodyFilename, r.Id)
		}
	}
	return nil
}

func (r *Response) HasBody() bool {
	return r.BodyFilename != nil
}

func GetResponseById(rs []Response, id int) (*Response, bool) {
	for _, r := range rs {
		if r.Id == id {
			return &r, true
		}
	}
	return nil, false
}