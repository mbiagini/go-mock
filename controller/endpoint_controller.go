package controller

import (
	"go-mock/apierrors"
	"go-mock/db"
	"go-mock/model"
	"go-mock/service"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/mbiagini/go-server-utils/gsrender"
	"github.com/mbiagini/go-server-utils/gsvalidation"
)

func GetEndpoints(w http.ResponseWriter, r *http.Request) {
	endpoints := db.DB.FindAll()
	gsrender.WriteJSON(w, http.StatusOK, endpoints)
}

func GetEndpointById(w http.ResponseWriter, r *http.Request) {
	strId := chi.URLParam(r, "id")
	err := gsvalidation.Var(strId, "number")
	if err != nil {
		gsrender.WriteJSON(w, http.StatusBadRequest, apierrors.NewWithMsg(apierrors.INVALID_ARGUMENT, err.Error()))
		return
	}
	id, _ := strconv.Atoi(strId)
	endpoint, ok := db.DB.FindById(id)
	if !ok {
		gsrender.WriteJSON(w, http.StatusNotFound, apierrors.New(apierrors.ENDPOINT_NOT_FOUND))
		return
	}
	gsrender.WriteJSON(w, http.StatusOK, endpoint)
}

func PostEndpoint(w http.ResponseWriter, r *http.Request) {
	var endpoint model.Endpoint
	if httpSuggestion := gsvalidation.DecodeJSONRequestBody(r, &endpoint); httpSuggestion != nil {
		code := apierrors.INTERNAL_SERVER_ERROR
		if httpSuggestion.Status != 500 {
			code = apierrors.INVALID_ARGUMENT
		}
		gsrender.WriteJSON(w, httpSuggestion.Status, apierrors.NewWithMsg(code, httpSuggestion.Message))
		return
	}
	if err := service.ValidateEndpoint(&endpoint); err != nil {
		gsrender.WriteJSON(w, http.StatusBadRequest, apierrors.NewWithMsg(apierrors.INVALID_ARGUMENT, err.Error()))
		return
	}
	endpoint, err := db.DB.Save(endpoint)
	if err != nil {
		gsrender.WriteJSON(w, http.StatusInternalServerError, apierrors.New(apierrors.INTERNAL_SERVER_ERROR))
		return
	}
	gsrender.WriteJSON(w, http.StatusCreated, endpoint)
}

func DeleteEndpoint(w http.ResponseWriter, r *http.Request) {
	strId := chi.URLParam(r, "id")
	err := gsvalidation.Var(strId, "number")
	if err != nil {
		gsrender.WriteJSON(w, http.StatusBadRequest, apierrors.NewWithMsg(apierrors.INVALID_ARGUMENT, err.Error()))
		return
	}
	id, _ := strconv.Atoi(strId)
	db.DB.Delete(id)
	gsrender.Status(w, http.StatusNoContent)
}

func UpdateEndpoint(w http.ResponseWriter, r *http.Request) {

	strId := chi.URLParam(r, "id")
	err := gsvalidation.Var(strId, "number")
	if err != nil {
		gsrender.WriteJSON(w, http.StatusBadRequest, apierrors.NewWithMsg(apierrors.INVALID_ARGUMENT, err.Error()))
		return
	}
	id, _ := strconv.Atoi(strId)
	_, ok := db.DB.FindById(id)
	if !ok {
		gsrender.WriteJSON(w, http.StatusNotFound, apierrors.New(apierrors.ENDPOINT_NOT_FOUND))
		return
	}

	var endpoint model.Endpoint
	if httpSuggestion := gsvalidation.DecodeJSONRequestBody(r, &endpoint); httpSuggestion != nil {
		code := apierrors.INTERNAL_SERVER_ERROR
		if httpSuggestion.Status != 500 {
			code = apierrors.INVALID_ARGUMENT
		}
		gsrender.WriteJSON(w, httpSuggestion.Status, apierrors.NewWithMsg(code, httpSuggestion.Message))
		return
	}

	if err := service.ValidateEndpoint(&endpoint); err != nil {
		gsrender.WriteJSON(w, http.StatusBadRequest, apierrors.NewWithMsg(apierrors.INVALID_ARGUMENT, err.Error()))
		return
	}

	endpoint.Id = id
	endpoint, err = db.DB.Save(endpoint)
	if err != nil {
		gsrender.WriteJSON(w, http.StatusInternalServerError, apierrors.New(apierrors.INTERNAL_SERVER_ERROR))
		return
	}
	gsrender.WriteJSON(w, http.StatusCreated, endpoint)
}