package service

import (
	"fmt"
	"go-mock/model"
	"net/http"

	"github.com/mbiagini/go-server-utils/gslog"
)

// ValidateEndpoint validates a given Endpoint struct. First checks whether its
// default response id corresponds to an existing response. Then, calls
// ValidateDiscriminator for any discriminator found (if any).
func ValidateEndpoint(e *model.Endpoint) error {
	if _, ok := GetResponseById(e.Responses, *e.DefaultResponseId); !ok {
		return fmt.Errorf("endpoint %s has a non-existing default response id %d", e.ToString(), *e.DefaultResponseId)
	}
	if len(e.Discriminators) != 0 {
		for _, d := range e.Discriminators {
			if err := ValidateDiscriminator(&d, e.Responses); err != nil {
				return fmt.Errorf("error in endpoint %s: %s", e.ToString(), err.Error())
			}
		}
	}
	return nil
}

// GetEndpointDefaultResponse returns its default Response, if found, or an error
// if its DefaultResponseId points to an unexisting response.
//
// Note: a valid endpoint should not be saved to DB if it has an invalid value for
// the field DefaultResponseId, so this function should never return an error.
func GetEndpointDefaultResponse(e model.Endpoint) (model.Response, error) {
	zero := model.Response{}
	resp, ok := GetResponseById(e.Responses, *e.DefaultResponseId)
	if !ok {
		return zero, fmt.Errorf("could not find default response for endpoint %s", e.ToString())
	}
	return resp, nil
}

// FindEndpointResponseFromRequest receives an Endpoint and an http.Request and
// returns the appropiate response to it. This response is found trying to match
// the request to any configured Discriminators. If non found, the default Response
// is returned.
//
// Note: this function can return an error if a Response cannot be found with an id.
// This means that the Endpoint was incorrectly configured and poorly validated, so
// it shouldn't happen.
func FindEndpointResponseFromRequest(e model.Endpoint, r *http.Request) (model.Response, error) {
	if len(e.Discriminators) == 0 {
		return GetEndpointDefaultResponse(e)
	}

	// Default response to return if no other applies.
	respId := *e.DefaultResponseId

	// Try to match any discriminator.
	for _, d := range e.Discriminators {

		id, ok, err := MatchDiscriminatorResponse(d, r)
		if err != nil {
			gslog.Server(fmt.Sprintf("Error analyzing request for endpoint %s: %s", e.ToString(), err.Error()))
		}

		// End loop when a discriminator's condition matches.
		if ok {
			respId = id
			break
		}
	}

	// Return response.
	zero := model.Response{}
	resp, ok := GetResponseById(e.Responses, respId)
	if !ok {
		return zero, fmt.Errorf("matched endpoint %s with responseId %d but no response found for it", e.ToString(), respId)
	}

	return resp, nil
}