package service

import (
	"fmt"
	"go-mock/model"
	"regexp"
	"strconv"
)

// ValidateCondition checks whether a given condition is well defined. Since a
// Condition has a regex and a responseId, this function first validates the regex
// and then searches for a Response in the received list that matches with the id
// found in the Condition.
func ValidateCondition(c *model.Condition, responses []model.Response) error {
	if _, err := regexp.Compile(c.Regexp); err != nil {
		return fmt.Errorf("regular expression %s found invalid: %s", c.Regexp, err.Error())
	}
	if _, ok := GetResponseById(responses, *c.ResponseId); !ok {
		return fmt.Errorf("couldn't find response for condition %s", c.ToString())
	}
	return nil
}

// MatchConditionFromValue receives a list of conditions and a string value and
// tries to match that value with any condition's regex. If any of those attempts
// returns true, that condition is returned.
//
// Note that more than one condition could match against the received value. This
// function will return the first condition that matches.
func MatchConditionFromValue(conditions []model.Condition, value string) (model.Condition, bool) {

	zero := model.Condition{}
	for _, c := range conditions {
		regex := regexp.MustCompile(c.Regexp)
		if regex.MatchString(value) {
			return c, true
		}
	}
	return zero, false
}

// MatchConditionFromStruct receives a list of conditions, a struct and a field name.
// It will search the struct for a field that matches the received one and, if found,
// will try to match its value against any of the given conditions, returning the
// first match, if found.
func MatchConditionFromStruct(conditions []model.Condition, v any, field string) (model.Condition, bool) {

	zero := model.Condition{}

	switch d := v.(type) {

	// struct.
	case map[string]interface{}:

		// iterate every struct's field.
		for name, value := range d {

			// found the field we are looking for.
			if name == field {
				// check if the value matches any of the conditions regexp.
				c, ok := matchFromField(conditions, value)
				if ok {
					return c, true
				}
			}

			// try to use value as struct in recursion.
			c, ok := MatchConditionFromStruct(conditions, value, field)
			if ok {
				return c, true
			}
		}

	// slice.
	case []interface{}:

		// iterate every element of the slice.
		for _, value := range d {

			// try to use value as struct in recursion.
			c, ok := MatchConditionFromStruct(conditions, value, field)
			if ok {
				return c, true
			}

		}

	}

	return zero, false

} 

// matchFromField searches in the given condition list if any of them matches against
// the specified field's value.
func matchFromField(conditions []model.Condition, v any) (model.Condition, bool) {
	zero := model.Condition{}
	for _, c := range conditions {
		regex := regexp.MustCompile(c.Regexp)
		match := fieldMatchesRegex(v, *regex)
		if match {
			return c, true
		}
	}
	return zero, false
}

// filedMatchesRegex receives a struct's field of undefined type and tries to match
// it against a regex. For this, reflection is used to determine the field's type.
func fieldMatchesRegex(v any, re regexp.Regexp) (bool) {
	switch v := v.(type) {
	case int:
		return re.MatchString(strconv.Itoa(v))
	case float64:
		return re.MatchString(strconv.FormatFloat(v, 'f', -1, 64))
	case string:
		return re.MatchString(v)
	default:
		return false
	}
}