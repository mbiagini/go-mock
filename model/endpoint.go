package model

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/mbiagini/go-server-utils/gslog"
)

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

type Standard string

const (
	REST Standard = "REST"
	SOAP Standard = "SOAP"
)

type Endpoint struct {
	Method              HttpMethod      `json:"method"`
	Path                string          `json:"path"`
	Standard            *Standard       `json:"standard"`
	DefaultResponseId   int             `json:"default_response_id"`	
	Discriminators      []Discriminator `json:"discriminators"`
	Responses           []Response      `json:"responses"`
}

func ValidateEndpoints(es []Endpoint) error {
	if len(es) == 0 {
		return errors.New("no endpoints received in configuration. A mock with no endpoints is not a mock to waste resources on")
	}
	for _, e := range es {
		err := e.Validate()
		if err != nil {
			return fmt.Errorf("configuration error: %s", err.Error())
		}
	}
	return nil
}

func (e *Endpoint) Validate() error {
	if e.Method == "" {
		return fmt.Errorf("invalid or missing method value for endpoint with path %s", e.Path)
	}
	if _, ok := GetResponseById(e.Responses, e.DefaultResponseId); !ok {
		return fmt.Errorf("endpoint %s has a non existing default response id configured: %d", e.ToString(), e.DefaultResponseId)
	}
	if len(e.Responses) == 0 {
		return fmt.Errorf("endpoint %s has no responses", e.ToString())
	}
	if len(e.Discriminators) != 0 {
		for _, d := range e.Discriminators {
			err := d.Validate(e.Responses, e.Standard)
			if err != nil {
				return fmt.Errorf("error in endpoint %s: %s", e.ToString(), err.Error())
			}
		}
	}
	for _, r := range e.Responses {
		err := r.Validate()
		if err != nil {
			return fmt.Errorf("error in endpoint %s: %s", e.ToString(), err.Error())
		}
	}
	return nil
}

func (e *Endpoint) GetDefaultResponse() (*Response, error) {
	resp, ok := GetResponseById(e.Responses, e.DefaultResponseId)
	if !ok {
		return nil, fmt.Errorf("could not find response for endpoint %s", e.ToString())
	}
	return resp, nil
}

func (e *Endpoint) FindResponse(r *http.Request) (*Response, error) {

	if len(e.Discriminators) == 0 {
		return e.GetDefaultResponse()
	}

	respId := e.DefaultResponseId

	for _, d := range e.Discriminators {

		id, err := d.FindResponseId(r)
		if err != nil {
			gslog.Server(fmt.Sprintf("error analyzing request for endpoint %s: %s", e.ToString(), err.Error()))
		}

		// end loop when discriminator condition matches.
		if id != nil {
			respId = *id
			break
		}

	}

	resp, ok := GetResponseById(e.Responses, respId)
	if !ok {
		return nil, fmt.Errorf("matched endpoint %s with responseId %d but no response found for it", e.ToString(), respId)
	}

	return resp, nil
}

func (e *Endpoint) ToString() string {
	return fmt.Sprintf("[method: %s, path: %s]", e.Method, e.Path)
}