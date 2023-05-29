package service

import "go-mock/model"

// GetResponseById receives a list of model.Response and an id. Searches the
// list for a Response that matches the given id and returns it.
func GetResponseById(rs []model.Response, id int) (model.Response, bool) {
	zero := model.Response{}
	for _, r := range rs {
		if *r.Id == id {
			return r, true
		}
	}
	return zero, false
}