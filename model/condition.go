package model

import (
	"fmt"
	"regexp"
	"strconv"
)

type Condition struct {
	Regexp     string `json:"regexp"`
	ResponseId int    `json:"response_id"`
}

func (c *Condition) Validate(rs []Response) error {
	if _, ok := GetResponseById(rs, c.ResponseId); !ok {
		return fmt.Errorf("couldn't find response for condition %s", c.ToString())
	}
	return nil
}

func (c *Condition) ToString() string {
	return fmt.Sprintf("[regexp: %s, response_id: %d]", c.Regexp, c.ResponseId)
}

func GetMatchingCondition(cs []Condition, value string) (*Condition, bool) {
	for _, c := range cs {
		re := regexp.MustCompile(c.Regexp)
		if re.MatchString(value) {
			return &c, true
		}
	}
	return nil, false
}

func GetMatchingConditionFromStruct(cs []Condition, v any, field string) (*Condition, error) {

	switch d := v.(type) {

	// struct.
	case map[string]interface{}:

		// iterate every struct's field.
		for name, value := range d {

			// found the field we are looking for.
			if name == field {

				// check if the value matches any of the conditions regexp.
				c, err := matchAny(cs, value)
				if err != nil {
					return nil, fmt.Errorf("found field %s but couldn't find matching condition: %s", field, err.Error()) 
				}
				if c != nil {
					return c, nil
				}

			}

			// try to use value as struct in recursion.
			c, err := GetMatchingConditionFromStruct(cs, value, field)
			if err != nil || c != nil {
				return c, err
			}

		}

	// slice.
	case []interface{}:

		// iterate every element of the slice.
		for _, value := range d {

			// try to use value as struct in recursion.
			c, err := GetMatchingConditionFromStruct(cs, value, field)
			if err != nil || c != nil {
				return c, err
			}

		}

	}

	return nil, nil
}

func matchAny(cs []Condition, v any) (*Condition, error) {
	for _, c := range cs {
		match, err := isMatch(v, c.Regexp)
		if err != nil {
			return nil, fmt.Errorf("could not compare with condition %s: %s", c.ToString(), err.Error())
		}
		if match {
			return &c, nil
		}
	}
	return nil, nil
}

func isMatch(v any, pattern string) (bool, error) {
	re := regexp.MustCompile(pattern)
	switch v := v.(type) {
	case int:
		return re.MatchString(strconv.Itoa(v)), nil
	case float64:
		return re.MatchString(strconv.FormatFloat(v, 'f', -1, 64)), nil
	case string:
		return re.MatchString(v), nil
	default:
		return false, fmt.Errorf("value has an unexpected type")
	}

}