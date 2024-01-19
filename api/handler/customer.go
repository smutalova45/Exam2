package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"main.go/api/models"
)

func (h Handler) Customer(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.CreateCustomer(w, r)
	case http.MethodGet:
		values := r.URL.Query()
		_, ok := values["id"]
		if !ok {
			h.GetCustomerList(w, r)
		} else {
			h.GetCustomerByID(w, r)
		}
	case http.MethodPut:
		h.UpdateCustomer(w, r)
	case http.MethodDelete:
		h.DeleteCustomer(w, r)
	}
}

func (h Handler) CreateCustomer(w http.ResponseWriter, r *http.Request) {
	createcustomer := models.CreateCustomer{}
	if err := json.NewDecoder(r.Body).Decode(&createcustomer); err != nil {
		handleResponse(w, http.StatusBadRequest, err)
		return
	}
	pk, err := h.storage.Customer().Create(createcustomer)
	if err != nil {
		handleResponse(w, 500, err)
		return
	}
	customer, err := h.storage.Customer().Get(models.PrimaryKey{
		ID: pk,
	})

	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}
	handleResponse(w, http.StatusCreated, customer)
}

func (h Handler) GetCustomerByID(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	v := values.Get("id")
	if len(v) <= 0 {
		handleResponse(w, http.StatusBadRequest, errors.New("id is required"))
		return
	}
	customer, err := h.storage.Customer().Get(models.PrimaryKey{ID: v})
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	handleResponse(w, http.StatusOK, customer)
}

func (h Handler) GetCustomerList(w http.ResponseWriter, r *http.Request) {
	var (
		page, limit = 1, 10
		search      string
		err         error
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
	resp, err := h.storage.Customer().GetList(models.GetListRequest{
		Page:   page,
		Limit:  limit,
		Search: search,
	})
	if err != nil {
		handleResponse(w, 500, err)
		return
	}
	handleResponse(w, 200, resp)

}

func (h Handler) UpdateCustomer(w http.ResponseWriter, r *http.Request) {
	updatecustomer := models.UpdateCustomer{}
	if err := json.NewDecoder(r.Body).Decode(&updatecustomer); err != nil {
		handleResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	pk, err := h.storage.Customer().Update(updatecustomer)
	if err != nil {
		handleResponse(w, 500, err)
		return
	}
	customer, err := h.storage.Customer().Get(models.PrimaryKey{
		ID: pk,
	})
	if err != nil {
		handleResponse(w, 500, err)
		return
	}

	handleResponse(w, 200, customer)
}

func (h Handler) DeleteCustomer(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	id := values.Get("id")
	if len(id) <= 0 {
		handleResponse(w, http.StatusBadRequest, errors.New("id is required"))
		return
	}
	if err := h.storage.Customer().Delete(models.PrimaryKey{
		ID: id,
	}); err != nil {
		handleResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	handleResponse(w, 200, "deleted data")
}
