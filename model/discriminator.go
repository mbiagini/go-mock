package model

import "fmt"

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

func (d *Discriminator) ToString() string {
	return fmt.Sprintf("[parameter: %s, location: %s, len(conditions): %d]", d.Parameter, d.Location, len(d.Conditions))
}