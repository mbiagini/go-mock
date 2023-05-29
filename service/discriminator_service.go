package service

import (
	"encoding/json"
	"fmt"
	"go-mock/model"
	"io/ioutil"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// ValidateDiscriminator calls ValidateCondition for every Condition found in the given
// Discriminator. If any error is found, it returns it.
func ValidateDiscriminator(d *model.Discriminator, responses []model.Response) error {
	for _, c := range d.Conditions {
		if err := ValidateCondition(&c, responses); err != nil {
			return err
		}
	}
	return nil
}

// MatchDiscriminatorResponse tries to find the given discriminator in an http.Request.
// If the discriminator's parameter is not found in the request, returns false. Unless
// there's an error. This can happen if the discriminator is found but then no condition
// matches the value found in the Request.
//
// If the discriminator is found and the Request's parameter value matches one of the
// conditions associated with the discriminator, returns the response's id.
func MatchDiscriminatorResponse(d model.Discriminator, r *http.Request) (int, bool, error) {

	zero := 0

	// Handle BODY discriminator.
	if d.Location == model.BODY {
		dataBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return zero, false, fmt.Errorf("error reading request's body: %s", err.Error())
		}
		var data interface{}
		if err := json.Unmarshal(dataBytes, &data); err != nil {
			return zero, false, fmt.Errorf("error unmarshaling request's body: %s", err.Error())
		}
		c, ok := MatchConditionFromStruct(d.Conditions, data, d.Parameter)
		if ok {
			return *c.ResponseId, true, nil
		}
		return zero, false, nil
	}

	value := ""

	switch d.Location {
	case model.PATH:
		value = chi.URLParam(r, d.Parameter)
	case model.QUERY:
		value = r.URL.Query().Get(d.Parameter)
	case model.HEADER:
		value = r.Header.Get(d.Parameter)
	}

	// Handle not found.
	if value == "" {
		return zero, false, nil
	}

	// Handle PATH / QUERY / HEADER discriminator.
	c, ok := MatchConditionFromValue(d.Conditions, value)
	if !ok {
		return zero, false, fmt.Errorf("parameter %s found in location %s but no matching condition found", d.Parameter, d.Location)
	}

	return *c.ResponseId, true, nil
}