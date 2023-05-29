package model

import (
	"fmt"
)

type Condition struct {
	Regexp     string `json:"regexp"      validate:"required"`
	ResponseId *int   `json:"response_id" validate:"required,number"`
}

func (c *Condition) ToString() string {
	return fmt.Sprintf("[regexp: %s, response_id: %d]", c.Regexp, c.ResponseId)
}