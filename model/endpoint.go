package model

import (
	"fmt"
)

// HttpMethod represents an enum used to define a method from HTTP protocol.
type HttpMethod string

const (
	GET     HttpMethod = "GET"
	POST    HttpMethod = "POST"
	PUT     HttpMethod = "PUT"
	DELETE  HttpMethod = "DELETE"
	PATCH   HttpMethod = "PATCH"
	HEAD    HttpMethod = "HEAD"
	OPTIONS HttpMethod = "OPTIONS"
)

// Standard represents an enum used to define the API standard of an endpoint.
type Standard string

const (
	REST Standard = "REST"
	SOAP Standard = "SOAP"
)

type Endpoint struct {
	Id                  int             `json:"id"`
	Method              HttpMethod      `json:"method"                   validate:"required,oneof=GET POST PUT DELETE PATCH HEAD OPTIONS"`
	Path                string          `json:"path"                     validate:"required"`
	Standard            Standard        `json:"standard"                 validate:"required,oneof=REST SOAP"`
	DefaultResponseId   *int            `json:"default_response_id"      validate:"required,number"`	
	Discriminators      []Discriminator `json:"discriminators,omitempty"`
	Responses           []Response      `json:"responses"                validate:"required,min=1,dive"` 
}

func (e *Endpoint) ToString() string {
	return fmt.Sprintf("[method: %s, path: %s]", e.Method, e.Path)
}