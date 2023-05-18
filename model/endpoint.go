package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/mbiagini/go-server-utils/gslog"

	"github.com/go-chi/chi/v5"
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
	Method              HttpMethod     `json:"method"`
	Path                string         `json:"path"`
	Standard            *Standard      `json:"standard"`
	HasDiscriminator    bool           `json:"has_discriminator"`
	DefaultResponseId   int            `json:"default_response_id"`	
	Discriminator       *Discriminator `json:"discriminator"`
	Responses           []Response     `json:"responses"`
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
	if e.HasDiscriminator {
		if e.Discriminator == nil {
			return fmt.Errorf("endpoint %s specifies it uses a discriminator but none found", e.ToString())
		}
		err := e.Discriminator.Validate(e.Responses, e.Standard)
		if err != nil {
			return fmt.Errorf("error in endpoint %s: %s", e.ToString(), err.Error())
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

	if !e.HasDiscriminator {
		return e.GetDefaultResponse()
	}

	d := e.Discriminator

	pValue := ""
	switch d.Location {
	case PATH:
		pValue = chi.URLParam(r, d.Parameter)
	case QUERY:
		pValue = r.URL.Query().Get(d.Parameter)
	case HEADER:
		pValue = r.Header.Get(d.Parameter)
	case BODY:
		dataBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return nil, fmt.Errorf("error reading request's body for endpoint %s: %s", e.ToString(), err.Error())
		}
		var data interface{}
		if err := json.Unmarshal(dataBytes, &data); err != nil {
			return nil, fmt.Errorf("error unmarshaling request's body for endpoint %s: %s", e.ToString(), err.Error())
		}
		c, err := GetMatchingConditionFromStruct(d.Conditions, data, d.Parameter)
		if err != nil {
			return nil, fmt.Errorf("error matching conditions for endpoint %s: %s", e.ToString(), err.Error())
		}
		if c == nil {
			gslog.Server(fmt.Sprintf("no matching condition found for parameter %s in endpoint %s", d.Parameter, e.ToString()))
			return e.GetDefaultResponse()
		}
		resp, ok := GetResponseById(e.Responses, c.ResponseId)
		if !ok {
			return nil, fmt.Errorf("matched endpoint %s with condition %s but no response found for it", e.ToString(), c.ToString())
		}
		return resp, nil

	}

	if pValue == "" {
		return nil, fmt.Errorf("discriminator parameter %s not found in location %s", d.Parameter, d.Location)
	}

	c, ok := GetMatchingCondition(d.Conditions, pValue)
	if !ok {
		gslog.Server(fmt.Sprintf("no matching condition found for parameter value %s in endpoint %s", pValue, e.ToString()))
		return e.GetDefaultResponse()
	}

	resp, ok := GetResponseById(e.Responses, c.ResponseId)
	if !ok {
		return nil, fmt.Errorf("matched endpoint %s with condition %s but no response found for it", e.ToString(), c.ToString())
	}

	return resp, nil
}

func (e *Endpoint) ToString() string {
	return fmt.Sprintf("[method: %s, path: %s]", e.Method, e.Path)
}