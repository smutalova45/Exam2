package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"main.go/api/models"
)

func (h Handler) Trip(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.CreateTrip(w, r)
	case http.MethodGet:
		values := r.URL.Query()
		if _, ok := values["id"]; !ok {
			h.GetTripList(w, r)
		} else {
			h.GetTripByID(w, r)
		}
	case http.MethodPut:
		h.UpdateTrip(w, r)
	case http.MethodDelete:
		h.DeleteTrip(w, r)
	}
}

func (h Handler) CreateTrip(w http.ResponseWriter, r *http.Request) {
	createtrip := models.CreateTrip{}
	if err := json.NewDecoder(r.Body).Decode(&createtrip); err != nil {
		handleResponse(w, http.StatusBadRequest, err)
		return
	}
	pk, err := h.storage.Trip().Create(createtrip)
	if err != nil {
		handleResponse(w, 500, err)
		return
	}
	trip, err := h.storage.Trip().Get(models.PrimaryKey{
		ID: pk,
	})

	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}
	handleResponse(w, http.StatusCreated, trip)
}

func (h Handler) GetTripByID(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	v := values.Get("id")
	if len(v) <= 0 {
		handleResponse(w, http.StatusBadRequest, errors.New("id is required"))
		return
	}
	trip, err := h.storage.Trip().Get(models.PrimaryKey{ID: v})
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	handleResponse(w, http.StatusOK, trip)
}

func (h Handler) GetTripList(w http.ResponseWriter, r *http.Request) {
	var (
		page, limit = 1, 10
		search      string

		err error
	)
	values := r.URL.Query()
	v := values.Get("page")
	if len(v) > 0 {
		page, err = strconv.Atoi(v)
		if err != nil {
			page = 1
		}
	}
	v1 := values.Get("limit")
	if len(v1) > 0 {
		limit, err = strconv.Atoi(v1)
		if err != nil {
			fmt.Println("limit", v1)
		}
	}
	v2 := values.Get("search")
	if len(v2) > 0 {
		search = v2
	}
	resp, err := h.storage.Trip().GetList(models.GetListRequest{
		Page:   page,
		Limit:  limit,
		Search: search,
	})
	if err != nil {
		fmt.Println("error getting list of trips", err.Error())
		handleResponse(w, 500, err)
		return
	}
	handleResponse(w, 200, resp)

}

func (h Handler) UpdateTrip(w http.ResponseWriter, r *http.Request) {
	updatetrip := models.Trip{}
	if err := json.NewDecoder(r.Body).Decode(&updatetrip); err != nil {
		handleResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	pk, err := h.storage.Trip().Update(updatetrip)
	if err != nil {
		handleResponse(w, 500, err)
		return
	}
	t, err := h.storage.Trip().Get(models.PrimaryKey{
		ID: pk,
	})
	if err != nil {
		handleResponse(w, 500, err)
		return
	}

	handleResponse(w, 200, t)
}

func (h Handler) DeleteTrip(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	id := values.Get("id")
	if len(id) <= 0 {
		handleResponse(w, http.StatusBadRequest, errors.New("id is required"))
		return
	}
	if err := h.storage.Trip().Delete(models.PrimaryKey{
		ID: id,
	}); err != nil {
		handleResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	handleResponse(w, 200, "deleted data")
}
