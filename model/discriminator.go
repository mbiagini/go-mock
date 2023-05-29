package model

import (
	"fmt"
)

type Location string

const (
	PATH   Location = "PATH"
	QUERY  Location = "QUERY"
	HEADER Location = "HEADER"
	BODY   Location = "BODY"
)

type Discriminator struct {
	Parameter  string      `json:"parameter"  validate:"required,alphanum"`
	Location   Location    `json:"location"   validate:"required,oneof=PATH QUERY HEADER BODY"`
	Conditions []Condition `json:"conditions" validate:"required,min=1"`
}

func (d *Discriminator) ToString() string {
	return fmt.Sprintf("[parameter: %s, location: %s, len(conditions): %d]", d.Parameter, d.Location, len(d.Conditions))
}