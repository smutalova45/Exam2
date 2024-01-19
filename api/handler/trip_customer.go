package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"main.go/api/models"
)

func (h Handler) TripCustomer(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.CreateTripCustomer(w, r)
	case http.MethodGet:
		values := r.URL.Query()
		if _, ok := values["id"]; !ok {
			h.GetTripCustomerList(w, r)
		} else {
			h.GetTripCustomerByID(w, r)
		}
	case http.MethodPut:
		h.UpdateTripCustomer(w, r)
	case http.MethodDelete:
		h.DeleteTripCustomer(w, r)
	}
}

func (h Handler) CreateTripCustomer(w http.ResponseWriter, r *http.Request) {
	create := models.CreateTripCustomer{}
	if err := json.NewDecoder(r.Body).Decode(&create); err != nil {
		handleResponse(w, http.StatusBadRequest, err)
		return
	}
	pk, err := h.storage.TripCustomer().Create(create)
	if err != nil {
		handleResponse(w, 500, err)
		return
	}
	tripcustomer, err := h.storage.Customer().Get(models.PrimaryKey{
		ID: pk,
	})

	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}
	handleResponse(w, http.StatusCreated, tripcustomer)
}

func (h Handler) GetTripCustomerByID(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	v := values.Get("id")
	if len(v) <= 0 {
		handleResponse(w, http.StatusBadRequest, errors.New("id is required"))
		return
	}
	tripcustomer, err := h.storage.TripCustomer().Get(models.PrimaryKey{ID: v})
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	handleResponse(w, http.StatusOK, tripcustomer)
}

func (h Handler) GetTripCustomerList(w http.ResponseWriter, r *http.Request) {
	var (
		page, limit = 1, 10

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
	resp, err := h.storage.TripCustomer().GetList(models.GetListRequest{
		Page:  page,
		Limit: limit,
	})
	if err != nil {
		handleResponse(w, 500, err)
		return
	}
	handleResponse(w, 200, resp)
}

func (h Handler) UpdateTripCustomer(w http.ResponseWriter, r *http.Request) {
	updatetripcustomer := models.TripCustomer{}
	if err := json.NewDecoder(r.Body).Decode(&updatetripcustomer); err != nil {
		handleResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	pk, err := h.storage.TripCustomer().Update(updatetripcustomer)
	if err != nil {
		handleResponse(w, 500, err)
		return
	}
	tripcustomer, err := h.storage.TripCustomer().Get(models.PrimaryKey{
		ID: pk,
	})
	if err != nil {
		handleResponse(w, 500, err)
		return
	}

	handleResponse(w, 200, tripcustomer)

}

func (h Handler) DeleteTripCustomer(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	id := values.Get("id")
	if len(id) <= 0 {
		handleResponse(w, http.StatusBadRequest, errors.New("id is required"))
		return
	}
	if err := h.storage.TripCustomer().Delete(models.PrimaryKey{
		ID: id,
	}); err != nil {
		handleResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	handleResponse(w, 200, "deleted data")
}
