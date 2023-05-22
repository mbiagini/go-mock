package model

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Location string

const (
	PATH   Location = "PATH"
	QUERY  Location = "QUERY"
	HEADER Location = "HEADER"
	BODY   Location = "BODY"
)

type Discriminator struct {
	Parameter  string      `json:"parameter"`
	Location   Location    `json:"location"`
	Conditions []Condition `json:"conditions"`
}

func (d *Discriminator) Validate(rs []Response, s *Standard) error {
	if d.Location == "" {
		return fmt.Errorf("discriminator %s has an invalid location value", d.ToString())
	}
	if d.Location == BODY {
		if s == nil || *s == "" {
			return fmt.Errorf("discriminator %s has a location of BODY but no or invalid architectual Standard defined", d.ToString())
		}
	}
	if len(d.Conditions) == 0 {
		return fmt.Errorf("no conditions found for discriminator of parameter %s", d.Parameter)
	}
	for _, c := range d.Conditions {
		err := c.Validate(rs)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *Discriminator) FindResponseId(r *http.Request) (*int, error) {

	// request parameter value.
	rpv := ""
	switch d.Location {
	case PATH:
		rpv = chi.URLParam(r, d.Parameter)
	case QUERY:
		rpv = r.URL.Query().Get(d.Parameter)
	case BODY:
		respId, err := d.FindResponseIdInBody(r)
		if err != nil {
			return nil, fmt.Errorf("error finding response in body: %s", err.Error())
		}
		return respId, nil
	}

	// discriminator not found in request.
	if rpv == "" {
		return nil, nil
	}

	// try to match a configured condition.
	c, ok := GetMatchingCondition(d.Conditions, rpv)
	if !ok {
		return nil, fmt.Errorf("parameter %s found in location %s but no matching condition found", d.Parameter, d.Location)
	}

	return &c.ResponseId, nil
}

func (d *Discriminator) FindResponseIdInBody(r *http.Request) (*int, error) {

	dataBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading request's body: %s", err.Error())
	}

	var data interface{}
	if err := json.Unmarshal(dataBytes, &data); err != nil {
		return nil, fmt.Errorf("error umarshaling request's body: %s", err.Error())
	}

	condition, err := GetMatchingConditionFromStruct(d.Conditions, data, d.Parameter)
	if err != nil {
		return nil, fmt.Errorf("error matching conditions: %s", err.Error())
	}

	if condition != nil {
		return &condition.ResponseId, nil
	}

	return nil, nil
}

func (d *Discriminator) ToString() string {
	return fmt.Sprintf("[parameter: %s, location: %s, len(conditions): %d]", d.Parameter, d.Location, len(d.Conditions))
}